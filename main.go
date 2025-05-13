package main

import (
  "flag"
  "fmt"
  "os"
  "os/exec"
  "net/http"
  "time"
  "encoding/json"
  "strings"
)

var options struct{
  port string
  templatepath string
}

type TIdcard struct{
  Type uint
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
  if d.Recepient.FullName == "" {return "1"}else{return "0"}
}

func split01(data string, partcount uint, capacity uint)(s0 string,s1 string,s2 string,s3 string){
  var result [4]string;
  result[0] = data;
  return result[0],result[1],result[2],result[3]
}

func split02(data string, partcount uint, capacity uint)(s0 string,s1 string,s2 string,s3 string){
  var result [4]string;
  a := strings.SplitN(data," ",4);for i:=0;i<len(a);i++ {result[i] = a[i]}
  return result[0],result[1],result[2],result[3]
}

func _DayOfMonth(d time.Time)string{
  return fmt.Sprintf("%02d",d.Day());
}
func _Month(d time.Time)string{
  return fmt.Sprintf("%02d",d.Month());
}
func _Year(d time.Time)string{
  return fmt.Sprintf("%04d",d.Year());  
}

func createFields(d TData)(r map[string]string){
  r = make(map[string]string);
  r["Text1"] = d.Provider.INN;
  r["Text2"] = d.Provider.KPP;
  r["Text3"] = d.ReportYear;
  r["Text4"] = d.CertNumber;
  r["Text5.0"] = d.CorrectionNumber;
  r["Text5.1"] = fmt.Sprintf("%03d",d.PagesCount());
  r["Text6.0"], r["Text6.1"], r["Text6.2"], r["Text6.3"] = split01(d.Provider.FullName,4,40);
  r["Text7.0"], r["Text7.1"], r["Text7.2"], _            = split02(d.Payer.FullName,3,36);
  r["Text8"] = d.Payer.INN;
  if d.Payer.INN == "" {
    r["Text9.1.0"] = _DayOfMonth(d.Payer.Idcard.Date);
    r["Text10.1.0"] = _Month(d.Payer.Idcard.Date);
    r["Text11.1.0"] = _Year(d.Payer.Idcard.Date);
    r["Text12"] = fmt.Sprintf("%02d",d.Payer.Idcard.Type);
    r["Text13"] = d.Payer.Idcard.SerNum;
  }
  r["Text9.0"] = _DayOfMonth(d.Payer.BD);
  //r["Text9.1.0"] = _DayOfMonth(d.Payer.Idcard.Date);
  r["Text9.1.1"] = _DayOfMonth(d.SignDate);
  r["Text10.0"] = _Month(d.Payer.BD);
  //r["Text10.1.0"] = _Month(d.Payer.Idcard.Date);
  r["Text10.1.1"] = _Month(d.SignDate);
  r["Text11.0"] = _Year(d.Payer.BD);
  //r["Text11.1.0"] = _Year(d.Payer.Idcard.Date);
  r["Text11.1.1"] = _Year(d.SignDate);
  //r["Text12"] = fmt.Sprintf("%02d",d.Payer.Idcard.Type);
  //r["Text13"] = d.Payer.Idcard.SerNum;
  r["Text14"] = d.SinglePage();
  r["Text15.0"] = fmt.Sprintf("%013d",uint(d.Total1 / 100));
  r["Text15.1"] = fmt.Sprintf("%013d",uint(d.Total2 / 100));
  r["Text16.0"] = fmt.Sprintf("%02d",uint(d.Total1 % 100));
  r["Text16.1"] = fmt.Sprintf("%02d",uint(d.Total2 % 100));
  r["Text17.0"], r["Text17.1"], r["Text17.2"], _         = split02(d.Signer.FullName,3,20);
  r["Text18.0"], r["Text18.1"], r["Text18.2"], _         = split02(d.Recepient.FullName,3,36);
  r["Text20"] = d.Recepient.INN;
  if d.Recepient.INN == "" {
    r["Text21.1"] = _DayOfMonth(d.Recepient.Idcard.Date);
    r["Text22.1"] = _Month(d.Recepient.Idcard.Date);
    r["Text23.1"] = _Year(d.Recepient.Idcard.Date);
    r["Text25"] = fmt.Sprintf("%02d",d.Recepient.Idcard.Type);
    r["Text26"] = d.Recepient.Idcard.SerNum;
  }
  r["Text21.0"] = _DayOfMonth(d.Recepient.BD);
  //r["Text21.1"] = _DayOfMonth(d.Recepient.Idcard.Date);
  r["Text22.0"] = _Month(d.Recepient.BD);
  //r["Text22.1"] = _Month(d.Recepient.Idcard.Date);
  r["Text23.0"] = _Year(d.Recepient.BD);
  //r["Text23.1"] = _Year(d.Recepient.Idcard.Date);
  //r["Text25"] = fmt.Sprintf("%02d",d.Recepient.Idcard.Type);
  //r["Text26"] = d.Recepient.Idcard.SerNum;
  r["Text30"] = d.SecondPageSignDateFull.Format("02.01.2006");
  return r;
}

func createResult(data TData)(*[]byte, error){
  var result []byte;
  s := ""; for k,v := range(createFields(data)){
    s += fmt.Sprintf("<field name=%q><value>%s</value></field>\n",k,v)
  }
  s = fmt.Sprintf(`
<?xml version="1.0" encoding="UTF-8"?>
<xfdf xmlns="http://ns.adobe.com/xfdf/" xml:space="preserve">
<fields>

%s

</fields>
</xfdf>
  `,s);

  var(
    f_result string = "./result.pdf"
    f_data string = "./data.xfdf"
    f_blank string = options.templatepath
  )

  os.WriteFile(f_data,[]byte(s),0777);
  result = []byte(s); //return &result,nil;

  d,e := os.Stat(options.templatepath);
  if (os.IsNotExist(e) || d.IsDir()) {}

  e = os.Remove(f_result);if e != nil {fmt.Printf("ERROR: %s\n\n",e.Error())}

  cmd := exec.Command("pdftk", f_blank, "fill_form", f_data, "output", f_result, "need_appearances");
  //cmd := exec.Command("pdftk " + f_blank + " fill_form ./data.xfdf output ./result.pdf need_appearances");
  e = cmd.Run();if e != nil {fmt.Printf("ERROR: %s\n\n",e.Error())}

  result,e = os.ReadFile(f_result);if e != nil {fmt.Printf("ERROR: %s\n\n",e.Error())}

  //os.Remove(f_data);
  //os.Remove(f_result);

  return &result,e;
}

func doOnGetPDF(w http.ResponseWriter, r *http.Request){

  data := TData{Provider: &TPerson{}, Payer: &TPerson{Idcard: TIdcard{}}, Recepient: &TPerson{Idcard: TIdcard{}}, Signer: &TPerson{}}
  var j []byte = make([]byte,r.ContentLength);
  r.Body.Read(j);
  fmt.Printf("BODY: %s\n\n",string(j));
  data.SignDate = time.Now();
  data.SecondPageSignDateFull = data.SignDate;
  e := json.Unmarshal(j,&data);if e != nil {println(e.Error())}
  result, _ := createResult(data);
  if 1 != 1 {fmt.Printf("TEST:\n%v\n=======\n%v\n\n",data,string(*result))}

  w.Header().Add("Access-Control-Allow-Origin","*");

  w.Write(*result);
}

func main(){
  //fmt.Printf("options:\n%s\n%s\n\n",options.port,options.templatepath);
  http.HandleFunc("/test",doOnGetPDF);
  http.ListenAndServe(options.port, nil);
}

func init(){
  flag.StringVar(&options.port,"port",":8080","Listen on port XXXXX. Default 8080");
  flag.StringVar(&options.templatepath,"template","./blank.pdf","Path to forms template. Default ./blank.pdf");
  flag.Parse();
}
