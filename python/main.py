from fastapi import FastAPI, HTTPException

from game import Game

app = FastAPI()


@app.post('/ai_attack')
async def receive_json(data: dict):
    game = Game(data)
    game.doTurn()
    return game.actions
    # except Exception as e:
    #     print(e)
    #     raise HTTPException(status_code=500, detail=str(e))


@app.post('/ai_upgrade')
async def receive_json(data: dict):
    game = Game(data)
    game.doUpgrade()
    return game.actions
