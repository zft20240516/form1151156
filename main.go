package main

import (
  "flag"
  "fmt"
  "os"
  "net/http"
  "time"
)

var options struct{
  port string
  templatepath string
}

type TIdcard struct{
  Type string
  SerNum string
  Date time.Time
}

type TPerson struct{
  FullName,INN,KPP string
  Idcard TIdcard
  BD time.Time
}

type TData struct{
  Provider, Payer, Recepient, Signer *TPerson
  ReportYear, CertNumber, CorrectionNumber string
  SignDate, SecondPageSignDateFull time.Time
  Total1, Total2 uint
}
func (d *TData)PagesCount()(uint){
  if d.SinglePage() == "1" {return 1}else{return 2}
}

func (d *TData)SinglePage()string{
  if d.Recepient == nil {return "1"}else{return "0"}
}

func split01(data string, partcount uint, capacity uint)(s0 string,s1 string,s2 string,s3 string){
  var result [4]string;
  return result[0],result[1],result[2],result[3]
}

func split02(data string, partcount uint, capacity uint)(s0 string,s1 string,s2 string,s3 string){
  var result [4]string;
  return result[0],result[1],result[2],result[3]
}

func _DayOfMonth(d time.Time)string{
  return fmt.Sprintf("%2d",d.Day());
}
func _Month(d time.Time)string{
  return fmt.Sprintf("%2d",d.Month());
}
func _Year(d time.Time)string{
  return fmt.Sprintf("%4d",d.Year());  
}

func createFields(d TData)(r map[string]string){
  r["Text1"] = d.Provider.INN;
  r["Text2"] = d.Provider.KPP;
  r["Text3"] = d.ReportYear;
  r["Text4"] = d.CertNumber;
  r["Text5.0"] = d.CorrectionNumber;
  r["Text5.1"] = fmt.Sprintf("%d",d.PagesCount());
  r["Text6.0"], r["Text6.1"], r["Text6.2"], r["Text6.3"] = split01(d.Provider.FullName,4,40);
  r["Text7.0"], r["Text7.1"], r["Text7.2"], _            = split02(d.Payer.FullName,3,36);
  r["Text8"] = d.Payer.INN;
  r["Text9.0"] = _DayOfMonth(d.Payer.BD);
  r["Text9.1.0"] = _DayOfMonth(d.Payer.Idcard.Date);
  r["Text9.1.1"] = _DayOfMonth(d.SignDate);
  r["Text10.0"] = _Month(d.Payer.BD);
  r["Text10.1.0"] = _Month(d.Payer.Idcard.Date);
  r["Text10.1.1"] = _Month(d.SignDate);
  r["Text11.0"] = _Year(d.Payer.BD);
  r["Text11.1.0"] = _Year(d.Payer.Idcard.Date);
  r["Text11.1.1"] = _Year(d.SignDate);
  r["Text12"] = fmt.Sprintf("%2d",d.Payer.Idcard.Type);
  r["Text13"] = d.Payer.Idcard.SerNum;
  r["Text14"] = d.SinglePage();
  r["Text15.0"] = fmt.Sprintf("%d",uint(d.Total1 / 100));
  r["Text15.1"] = fmt.Sprintf("%d",uint(d.Total2 / 100));
  r["Text16.0"] = fmt.Sprintf("%d",uint(d.Total1 % 100));
  r["Text16.1"] = fmt.Sprintf("%d",uint(d.Total2 % 100));
  r["Text17.0"], r["Text17.1"], r["Text17.2"], _         = split02(d.Signer.FullName,3,20);
  r["Text18.0"], r["Text18.1"], r["Text18.2"], _         = split02(d.Recepient.FullName,3,36);
  r["Text20"] = d.Recepient.INN;
  r["Text21.0"] = _DayOfMonth(d.Recepient.BD);
  r["Text21.1"] = _DayOfMonth(d.Recepient.Idcard.Date);
  r["Text22.0"] = _Month(d.Recepient.BD);
  r["Text22.1"] = _Month(d.Recepient.Idcard.Date);
  r["Text23.0"] = _Year(d.Recepient.BD);
  r["Text23.1"] = _Year(d.Recepient.Idcard.Date);
  r["Text25"] = fmt.Sprintf("%2d",d.Recepient.Idcard.Type);
  r["Text26"] = d.Recepient.Idcard.SerNum;
  r["Text30"] = d.SecondPageSignDateFull.Format("02.01.2006");
  return r;
}

func createResult()(*[]byte, error){
  d,e := os.Stat(options.templatepath);
  if (os.IsNotExist(e) || d.IsDir()) {}
  result := []byte("TEST");
  return &result,nil;
}

func main(){
  fmt.Printf("options:\n%s\n%s\n\n",options.port,options.templatepath);
  http.ListenAndServe(":8080", nil);
}

func init(){
  flag.StringVar(&options.port,"port",":8080","Listen on port XXXXX. Default 8080");
  flag.StringVar(&options.templatepath,"template","./template.pdf","Path to forms template. Default ./template.pdf");
  flag.Parse();
}
