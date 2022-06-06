import typing

from fastapi import Response

from piccolo.utils.pydantic import create_pydantic_model
from piccolo.engine import engine_finder

from blacksheep import Application, FromJSON, json, status_code, bad_request, Response, Content

from sql_app.tables import Expense

app = Application()

ExpenseModelIn: typing.Any = create_pydantic_model(table=Expense, model_name=" ExpenseModelIn")

ExpenseModelOut: typing.Any = create_pydantic_model(table=Expense, include_default_columns=True, model_name=" ExpenseModelIn")

ExpenseModelPartial: typing.Any = create_pydantic_model(
    table=Expense, model_name="ExpenseModelPartial", all_optional=True
)

@app.router.get("/expenses")
async def expenses():
    try:
        expense = await Expense.select()
        return expense
    except:
        return Response(404, content=Content(b"text/plain", b"Not Found"))


@app.router.get("/expense/{id}")
async def expense(id: int):
    expense = await Expense.select().where(id==Expense.id)
    if not expense:
        return Response(404, content=Content(b"text/plain", b"Id not Found"))
    return expense

@app.router.post("/expense")
async def create_expense(expense_model: FromJSON[ExpenseModelIn]):
    
    try:
        expense = Expense(**expense_model.value.dict())
        await expense.save()
        return ExpenseModelOut(**expense.to_dict())  
    except:
        return Response(400, content=Content(b"text/plain", b"Bad Request"))

@app.router.patch("expense/{id}")
async def patch_expense(
        id: int, expense_model: FromJSON[ExpenseModelPartial]
):
    expense = await Expense.objects().get(Expense.id == id)
    if not expense:
        return Response(404, content=Content(b"text/plain", b"Id not Found"))

    for key, value in expense_model.value.dict().items():
        if value is not None:
            setattr(expense, key, value)

    await expense.save()
    return ExpenseModelOut(**expense.to_dict())


@app.router.put("/expense/{id}")
async def put_expense(
        id: int, expense_model: FromJSON[ExpenseModelIn]
):
    expense = await Expense.objects().get(Expense.id == id)
    if not expense:
        return Response(404, content=Content(b"text/plain", b"Id not Found"))

    for key, value in expense_model.value.dict().items():
        if value is not None:
            setattr(expense, key, value)

    await expense.save()
    return ExpenseModelOut(**expense.to_dict())

@app.router.delete("/expense/{id}")
async def delete_expense(id: int):
    expense = await Expense.objects().get(Expense.id == id)
    if not expense:
        return Response(404, content=Content(b"text/plain", b"Id Not Found"))
    await expense.remove()
    return json({"message":"Expense deleted"})



async def open_database_connection_pool(application):
    try:
        engine = engine_finder()
        await engine.start_connection_pool()
    except Exception:
        print("Unable to connect to the database")


async def close_database_connection_pool(application):
    try:
        engine = engine_finder()
        await engine.close_connection_pool()
    except Exception:
        print("Unable to connect to the database")


app.on_start += open_database_connection_pool
app.on_stop += close_database_connection_pool    