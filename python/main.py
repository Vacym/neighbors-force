from fastapi import FastAPI, HTTPException
# from time import time

from game import Game

app = FastAPI()
game = Game()


@app.post('/ai_attack')
async def receive_json(data: dict):
    game.reinit(data)
    game.doTurn()
    return game.actions


@app.post('/ai_upgrade')
async def receive_json(data: dict):
    game.reinit(data)
    game.doUpgrade()
    return game.actions
