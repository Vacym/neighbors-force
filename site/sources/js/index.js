function stringObjectToNumbers(formData) {
    const obj = {};
    Object.keys(formData).forEach((key, index) => {
        obj[key] = parseInt(formData[key]);
    });
    return obj;
}

function sendFormData(formData) {
    const numPlayers = parseInt(formData.get('num_players'));

    const botLevelInputs = new Array(numPlayers)
        .fill(0)
        .map((_, index) => {
            const inputName = `bot_levels[${index}]`;
            const inputValue = formData.get(inputName);
            return inputValue !== null ? parseInt(inputValue) : 0;
        });

    const requestData = {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            rows: parseInt(formData.get('rows')),
            cols: parseInt(formData.get('cols')),
            num_players: numPlayers,
            player_id: parseInt(formData.get('player_id')),
            bot_levels: botLevelInputs
        })
    };

    return fetch('api/game/create', requestData)
        .then(response => response.json());
}

async function fetchMap() {
    const response = await fetch('api/game/get_map');
    if (!response.ok) {
        throw new Error('Failed to fetch the map');
    }
    return response.json();
}

function sendData(data, path) {
    const requestData = {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data)
    };

    return fetch(path, requestData)
        .then(response => response.json());
}

function endAttack() {
    sendData({}, "api/game/end_attack").
        then(data => startUpgrade(data))
        .catch(error => console.error(error));
}

function endTurn() {
    sendData({}, "api/game/end_turn").
        then(data => startNewGame(data))
        .catch(error => console.error(error));
}


async function recover(error) {
    try {
        const data = await fetchMap();
        if (data.players[data.turn].attacking === true) {
            startNewGame(data);
        } else {
            startUpgrade(data);
        }
        if (error) renderError(error)
    } catch (error) {
        console.error('Error while restoring the game:', error);
    }
}

function startNewGame(data) {
    game = data
    const board = data.board;
    const players = data.players;
    const turn = data.turn

    renderError(data.error)
    if (data.error) {
        recover(data.error)
        return
    }

    boardElement = renderNewBoard(board)
    renderScores(players[turn].points)
    markCanAttack(boardElement, turn)
    addAttackClickHandlers(boardElement)
}

function startUpgrade(data) {
    game = data
    const board = data.board;
    const players = data.players;
    const turn = data.turn

    renderError(data.error)
    if (data.error) {
        recover(data.error)
        return
    }

    boardElement = renderNewBoard(board)
    renderScores(players[turn].points)
    markCanUpgrade(boardElement, turn)
    addUpgradeClickHandlers(boardElement)
}

function renderError(error) {
    errorDiv = document.getElementById('error');
    if (error) {
        errorDiv.textContent = error
    } else {
        errorDiv.textContent = "there is will be error"
    }
}

function renderScores(scores) {
    const scoreElement = document.getElementById('scores');
    scoreElement.textContent = scores;
}

function renderNewBoard(board) {
    const boardElement = document.getElementById("board");
    boardElement.innerHTML = "";

    for (let row = 0; row < board.rows; row++) {
        const table = document.createElement("table");
        table.id = "table-" + row;
        const tbody = document.createElement("tbody");
        table.appendChild(tbody);
        for (let col = 0; col < board.cols; col++) {
            const cell = board.cells[row][col];
            const td = document.createElement("td");
            td.id = row * board.cols + col;
            if (cell != undefined) {
                td.setAttribute("power", cell.power);
                if (cell.owner_id >= 0) {
                    td.setAttribute("owner-id", cell.owner_id);
                }

                const divPower = document.createElement("div");
                divPower.classList.add("num", "power");
                divPower.innerText = cell.power;

                const divLevel = document.createElement("div");
                divLevel.classList.add("num", "level");
                divLevel.innerText = cell.level;

                const svg = document.createElementNS("http://www.w3.org/2000/svg", "svg");

                const use = document.createElementNS("http://www.w3.org/2000/svg", "use");
                use.setAttribute("href", "#hexagon");
                svg.appendChild(use);

                td.appendChild(svg);
                td.appendChild(divPower);
                td.appendChild(divLevel);
            }
            tbody.appendChild(td);
        }
        boardElement.appendChild(table);
    }

    return boardElement
}

function markCanAttack(boardElement, turn) {
    const tables = boardElement.getElementsByTagName('table');
    for (let i = 0; i < tables.length; i++) {
        const cells = tables[i].getElementsByTagName('td');
        for (let j = 0; j < cells.length; j++) {
            const cell = cells[j];
            const owner_id = parseInt(cell.getAttribute('owner-id'));
            if (owner_id == turn) {
                cell.classList.add('can-attack');
            } else {
                cell.classList.remove('can-attack');
            }
        }
    }
}

function addAttackClickHandlers(boardElement) {
    const tables = boardElement.querySelectorAll('table');

    tables.forEach(table => {
        const tds = table.querySelectorAll('td');

        tds.forEach(td => {
            if (td.classList.contains('can-attack')) {
                td.onclick = (event) => { addCanBeAttackedClickHandlers(boardElement, event) }
            } else {
                td.onclick = null;
            }
        });
    });
};

function addCanBeAttackedClickHandlers(boardElement, event) {
    const tables = boardElement.querySelectorAll('table');
    let attackCellCoords = determineCoords(event.currentTarget.id, game.board.rows, game.board.cols)
    console.log("attackCellCoords", attackCellCoords)
    console.log('click attack!');

    markAttacking(boardElement, event.currentTarget)
    markCanBeAttacked(boardElement, this)

    tables.forEach(table => {
        const allTds = table.querySelectorAll('td');
        allTds.forEach(otherTd => {
            if (otherTd.classList.contains('can-be-attacked')) {
                otherTd.onclick = (event) => {
                    console.log('Other cell attacked!');
                    let attackedCellCoords = determineCoords(event.currentTarget.id, game.board.rows, game.board.cols)
                    console.log("attackedCellCoords", attackedCellCoords)
                    removeCanBeAttackedClickHandlers(boardElement)
                    unmarkCanBeAttacked(boardElement)

                    sendData({ from: attackCellCoords, to: attackedCellCoords }, '/api/game/attack')
                        .then(data => startNewGame(data))
                        .catch(error => console.error(error));
                };
            }
        })
    });
}

function removeCanBeAttackedClickHandlers(boardElement) {
    const tables = boardElement.querySelectorAll('table');

    tables.forEach(table => {
        const tds = table.querySelectorAll('td');

        tds.forEach(td => {
            if (td.classList.contains('can-be-attacked')) {
                td.onclick = null
            }
        })
    })
}


function markCanBeAttacked(boardElement, attackTd) {
    const tables = boardElement.querySelectorAll('table');

    tables.forEach(table => {
        const tds = table.querySelectorAll('td');

        tds.forEach(td => {
            if (!td.classList.contains('can-attack')) {
                td.classList.add('can-be-attacked')
            }
        })
    })
}

function unmarkCanBeAttacked(boardElement, attackTd) {
    const tables = boardElement.querySelectorAll('table');

    tables.forEach(table => {
        const tds = table.querySelectorAll('td');

        tds.forEach(td => {
            if (td.classList.contains('can-be-attacked')) {
                td.classList.remove('can-be-attacked')
            }
        })
    })
}

function markAttacking(boardElement, newAttacking) {
    if (attacking) {
        let lastAttacking = document.getElementById(attacking.id)
        lastAttacking.classList.remove('attacking')
    }

    attacking = newAttacking
    attacking.classList.add('attacking')
}


function markCanUpgrade(boardElement, turn) {
    const tables = boardElement.getElementsByTagName('table');
    for (let i = 0; i < tables.length; i++) {
        const cells = tables[i].getElementsByTagName('td');
        for (let j = 0; j < cells.length; j++) {
            const cell = cells[j];
            const owner_id = parseInt(cell.getAttribute('owner-id'));
            if (owner_id == turn) {
                cell.classList.add('can-upgrade');
            } else {
                cell.classList.remove('can-upgrade');
            }
        }
    }
}

function addUpgradeClickHandlers(boardElement) {
    const tables = boardElement.querySelectorAll('table');

    tables.forEach(table => {
        const tds = table.querySelectorAll('td');

        tds.forEach(td => {
            if (td.classList.contains('can-upgrade')) {
                td.onclick = (event) => {
                    console.log('Cell upgraded!');
                    let upgradedCellCoords = determineCoords(event.currentTarget.id, game.board.rows, game.board.cols)
                    console.log("upgradedCellCoords", upgradedCellCoords)

                    sendData({ cell: upgradedCellCoords, levels: 1 }, '/api/game/upgrade')
                        .then(data => startUpgrade(data))
                        .catch(error => console.error(error));
                };
            }
        });
    });
};

function determineCoords(id, rows, cols) {
    id = Number(id)
    return { row: Math.floor(id / rows), col: id % cols }
}

var game
var attacking

function main() {
    const form = document.querySelector('#game-form');
    form.addEventListener('submit', function (event) {
        event.preventDefault();

        const formData = new FormData(form);
        console.log(formData)
        sendFormData(formData)
            .then(data => startNewGame(data))
            .catch(error => console.error(error));
    });

    const endAttackButton = document.querySelector('#end-attack');
    endAttackButton.onclick = endAttack
    const endTurnButton = document.querySelector('#end-turn');
    endTurnButton.onclick = endTurn

    // Process keyboard events
    document.addEventListener('keydown', function (event) {
        if (event.key === 'N' || event.key === 'n') {
            form.dispatchEvent(new Event('submit'));
        } else if (event.key === 'Z' || event.key === 'z') {
            endAttackButton.click();
        } else if (event.key === 'X' || event.key === 'x') {
            endTurnButton.click();
        }
    });
}

document.addEventListener('DOMContentLoaded', main);


