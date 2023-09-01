from random import randint


class Game:
    def reinit(self, game: dict) -> None:
        self.player = game['turn']
        self.game = game['board']['cells']
        self.points = game['players'][self.player]['points']
        self.actions = {'attack': None, 'upgrade': None}

    def doTurn(self) -> None:
        best_score, best_from, best_to = 0, (0, 0), (0, 0)
        for row in range(len(self.game)):
            for col in range(len(self.game[row])):
                if self.game[row][col] is None:
                    continue

                cell = (row, col)
                cell_power = self.get_cell(cell)['power']
                if self.owner(cell) == self.player and cell_power > 1:
                    for to in self.get_alien_neighbors(cell):
                        score = self.calculate_score(cell_power, to)
                        if score > best_score:
                            best_score = score
                            best_from = cell
                            best_to = to

        # print(f'FINAL: {best_score=}, {best_from=}, {best_to=}')
        if best_score:
            self.put_attack(best_from, best_to)

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
            upgradeCost = cell_level * (cell_level + 1) // 2
            if self.points >= upgradeCost:
                return self.put_upgrade(cell)

    def get_alien_neighbors(self, cell: tuple) -> None:
        offset = cell[0] % 2
        neighborsRelative = (
            (-1, offset - 1),  # up-left
            (-1, offset - 0),  # up-right
            (+1, offset - 1),  # down-left
            (+1, offset - 0),  # down-right
            (0, -1),           # left
            (0, +1),           # right
        )

        for row, col in neighborsRelative:
            to = (cell[0] + row, cell[1] + col)
            if (self.get_cell(to) is not None
               and self.owner(to) != self.player):
                yield to

    def calculate_score(self, from_power, to) -> int:
        # If cell is free
        if self.owner(to) is None:
            return from_power

        # If our cell power lower than other
        to_power = self.get_cell(to)['power']
        if from_power <= to_power:
            return 100 + (to_power - from_power)

        # Otherwise
        return 10000 + (to_power - from_power)

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
