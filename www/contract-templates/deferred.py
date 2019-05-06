import smart_transaction
import time_date_engine


def complex_entry_point_0():
  smart_transaction.send('{{.Receiver}}', '{{.Data}}', '{{.TransactionMessage}}')

def start_rules():
  time_date_engine.start_on_date('{{.ContractDate}}', complex_entry_point_0)
