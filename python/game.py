from random import randint


class Game:
    def __init__(self, game: dict) -> None:
        self.game = game['board']['cells']
        self.player = game['turn']
        self.points = game['players'][self.player]['points']
        self.actions = {'attack': None, 'upgrade': None}

    def doTurn(self) -> None:
        for row in range(len(self.game)):
            for col in range(len(self.game[row])):
                if self.game[row][col] is None:
                    continue

                cell = (row, col)
                if (self.owner(cell) == self.player
                   and self.get_cell(cell)['power'] > 1):
                    self.attack(cell)

    def doUpgrade(self) -> None:
        if not self.points:
            return

        ownedCells = []
        for row in range(len(self.game)):
            for col in range(len(self.game[row])):
                if self.game[row][col] is None:
                    continue

                cell = (row, col)
                if self.owner(cell) == self.player:
                    ownedCells.append(cell)

        maxrow, maxcol = len(self.game), len(self.game[0])
        main_cell = {
            0: (0, 0),
            1: (maxrow, maxcol),
            2: (maxrow, 0),
            3: (0, maxcol)
        }[self.player]

        ownedCells.sort(key=lambda cell:
                        [abs(main_cell[0]-cell[0]) + abs(main_cell[1]-cell[1]),
                         -self.get_cell(cell)['level']],
                        reverse=True)

        for cell in ownedCells:
            cell_level = self.get_cell(cell)['level']
            upgradeCost = (cell_level * (cell_level + 1)) // 2
            if self.points >= upgradeCost:
                self.put_upgrade(cell)
                break

    def attack(self, cell: tuple) -> None:
        neighbors = self.get_alien_neighbors(cell)
        if len(neighbors):
            to = neighbors[randint(0, len(neighbors)-1)]
            self.put_attack(cell, to)

    def get_alien_neighbors(self, cell: tuple) -> list:
        offset = cell[0] % 2
        neighborsRelative = (
            (-1, offset - 1),  # up-left
            (-1, offset - 0),  # up-right
            (+1, offset - 1),  # down-left
            (+1, offset - 0),  # down-right
            (0, -1),           # left
            (0, +1),           # right
        )

        neighbors = []
        for row, col in neighborsRelative:
            to = (cell[0] + row, cell[1] + col)
            if (self.get_cell(to) is not None
               and self.owner(to) != self.player):
                neighbors.append(to)
        return neighbors

    def put_attack(self, cell: tuple, to: tuple) -> None:
        self.actions['attack'] = [cell, to]

    def put_upgrade(self, cell: tuple) -> None:
        self.actions['upgrade'] = list(cell)

    def get_cell(self, cell: tuple) -> dict:
        if (cell[0] < 0 or cell[1] < 0 or cell[0] >= len(self.game)
           or cell[1] >= len(self.game[cell[0]])):
            return None
        return self.game[cell[0]][cell[1]]

    def owner(self, cell: tuple) -> int:
        return self.get_cell(cell).get('owner_id', None)
