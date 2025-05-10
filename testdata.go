package main

var TEST_j []byte = []byte(`
{
  "Provider": {
    "INN": "740499999997",
    "FullName": "ИП Иванов Иван Иванович"
  },
  "Signer": {
    "FullName": "Иванов Иван Иванович"
  },
  "Payer": {
    "BD": "2001-02-27T00:00:00Z",
    "INN": "740499999999",
    "FullName": "Смирнова Елена Александровна"
  },
  "Recepient": {
    "BD": "2016-03-31T00:00:00Z",
    "-INN": "74049999998",
    "Idcard": {
      "SerNum": "001 9876543",
      "Date": "2021-04-30T00:00:00Z",
      "Type": 3
    },
    "-FullName": "Кузнецов Андрей Сергеевич"
  },
  "ReportYear": "2025",
  "CertNumber": "321",
  "CorrectionNumber": "2",
  "Total1": 987654321,
  "Total2": 9987654321 
}
  `);
