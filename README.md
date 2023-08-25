# Neighbors Force - Hexagonal Field Strategy Game

Welcome to Neighbors Force, an engaging computer strategy game written in Go that takes place on a hexagonal field.

## Game Overview

Neighbors Force is a game where players compete to dominate a hexagonal field by capturing cells and overpowering their opponents.

## Objective

The main objective of the game is to capture as many cells on the hexagonal field as possible while strategically eliminating opponents.

## Hexagonal Field

The game's field is composed of hexagonal cells, each featuring two essential attributes: energy and level. The energy of a cell determines its battle strength, while the level influences its capacity to enhance neighboring cells.

Players' cells are uniquely colored, with no gray cells on the field. At the start, each player possesses a single cell on the field.

## Gameplay Phases

The gameplay is divided into turns, each encompassing two phases: Attack and Upgrade.

### Attack Phase

In the Attack phase, players can attempt to capture unclaimed or opponent-owned neighboring cells. By selecting one of their cells, players can channel energy into it. If the attacking cell's energy surpasses the targeted cell's energy, the attacker captures the cell. If the attacker's energy is lower, the opposing cell's energy decreases to 1.

Players can execute multiple attacks in a turn, depending on the energy available on their cells.

### Upgrade Phase

During the Upgrade phase, players earn points based on the number of cells they've captured. These points can be used to enhance their cells. Subsequent upgrades cost more than the previous ones. Upgraded cells boost the energy of adjacent cells, excluding the upgraded cell itself.

Players can upgrade multiple cells in a turn, provided they have sufficient points.

## Endgame

The game continues until a single player remains or the predefined number of turns expires.

The winner is determined by the last surviving player or the one who has captured the largest territory.


## License

This game is released under the [MIT License](LICENSE).