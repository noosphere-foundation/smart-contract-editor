import smart_transaction
import time_date_engine
import smart_exchange

def complex_entry_point_0():
  smart_transaction.send('{{.Receiver}}', '{{.Data}}', '{{.TransactionMessage}}')

def start_rules():
  if smart_exchange.getER('{{.SelectCondition}}') {{.SelectOperator}} {{.SelectValue}}:
    complex_entry_point_0()

  time_date_engine.start_on_date('{{.ContractDate}}', complex_entry_point_0)
