from piccolo.table import Table
from piccolo.columns import Varchar, Integer

class Expense(Table):
    amount = Integer()
    description = Varchar()
