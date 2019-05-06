import smart_transaction
import time_date_engine


def complex_entry_point_0():
  for receiver in {{.Receivers}}:
    smart_transaction.send(receiver, '{{.Data}}', '{{.TransactionMessage}}')

def start_rules():
  smart_transaction.start_on_date('{{.ContractDate}}', complex_entry_point_0)
