import smart_transaction
import time_date_engine


def auto_payment_start(mode, complex_entry_point):
  if mode == 'every_day':
    smart_transaction.start_every_day(complex_entry_point)
  elif mode == 'every_week':
    smart_transaction.start_every_week(complex_entry_point)
  # what's about the leap year?
  elif mode == 'every_month':
    smart_transaction.start_every_month(complex_entry_point)
  elif mode == 'every_year':
    smart_transaction.start_every_year(complex_entry_point)

def complex_entry_point_0():
  smart_transaction.send('{{.Receiver}}', '{{.Data}}', '{{.TransactionMessage}}')

def start_rules():
  auto_payment_start('{{.AutoPaymentMode}}', complex_entry_point_0)
  smart_transaction.start_on_date('{{.ContractDate}}', complex_entry_point_0)
