package main

import (
	"database/sql"
	"log"
	"net/http"
	"fmt"
	"io/ioutil"
	"strconv"
	
        "html/template"
        _ "github.com/lib/pq"
        

	
)

var db *sql.DB


func getClient(clientID int) (Client, error) {
	//Retrieve
	res := Client{}

	var id int
	var hostname string
	var ip string
	var dt string
	var tm string

	err := db.QueryRow(`SELECT id, hostname, ip, dt, tm FROM clients where id = $1`, clientID).Scan(&id, &hostname, &ip, &dt, &tm)
	if err == nil {
		res = Client{ID: id, Hostname: hostname, IP: ip, Dt: dt, Tm: tm}
	}

	return res, err
}

func allClients() ([]Client, error) {
	//Retrieve
	clients := []Client{}

	rows, err := db.Query(`SELECT * FROM clients`)
	defer rows.Close()
	if err == nil {
		for rows.Next() {
			var id int
			var hostname string
			var ip string
			var dt string
			var tm string
                        
			err = rows.Scan(&id, &hostname, &ip, &dt, &tm)
		
			if err == nil {
				currentClient := Client{ID: id, Hostname: hostname, IP: ip, Dt: dt, Tm: tm}
				fmt.Println(currentClient)
				clients = append(clients, currentClient)
				
			} else {
				return clients, err
			}
		}
	} else {
		return clients, err
	}

	return clients, err
	
}


func allCerts()  ([]Cert, error){
	//Retrieve
	certs := []Cert{}

	rows, err := db.Query(`SELECT * FROM certs`)
	//fmt.Println(rows)
	//fmt.Fprintf("i am here")
	defer rows.Close()
	if err == nil {
		for rows.Next() {
					var certname string
			var selectionval string
			
                        
			err = rows.Scan(&certname, &selectionval)
		
			if err == nil {
				currentCert := Cert{Certname: certname, Selectionval: selectionval}
				//fmt.Println(currentCert)
				//fmt.Fprint(currentCert)
				certs = append(certs, currentCert)
				
			} else {
				return certs, err
			}
		}
	} else {
		return certs, err
	}

	return certs, err
	
	
}
//func insertcertselection(certname string, selectionval string) (int, error) {
	//Create
	
//	err1 := db.QueryRow(`INSERT INTO certs(certname, selectionval) VALUES($1, $2)`, certname, selectionval)
//	return 0,err1
//}INSERT INTO certs(certname, selectionval) VALUES($1, $2)

func inserting(certname string, productsSelected string) (int, error) {
	//Create
	var clientID int
	err := db.QueryRow(`INSERT INTO certs(certname, selectionval) VALUES($1, $2) `, certname, productsSelected).Scan(&clientID)
         if err != nil {
		return 0, err
	}

	fmt.Printf("Last inserted ID: %v\n", clientID)
	return 0, err
	//return clientID, err
}
func insertClient(hostname string, ip string, dt string, tm  string) (int, error) {
	//Create
	var clientID int
	err := db.QueryRow(`INSERT INTO clients(hostname, ip, dt, tm) VALUES($1, $2, $3, $4) RETURNING id`, hostname, ip, dt, tm).Scan(&clientID)

	if err != nil {
		return 0, err
	}

	fmt.Printf("Last inserted ID: %v\n", clientID)
	return clientID, err
}

func updateClient(id int, hostname string, ip string, dt string, tm string) (int, error) {
	//Create
	res, err := db.Exec(`UPDATE clients set hostname=$1, ip=$2, dt=$3, tm=$4 where id=$5 RETURNING id`, hostname, ip, dt, tm, id)
	if err != nil {
		return 0, err
	}

	rowsUpdated, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsUpdated), err
}

func removeClient(clientID int) (int, error) {
	//Delete
	res, err := db.Exec(`delete from clients where id = $1`, clientID)
	if err != nil {
		return 0, err
	}

	rowsDeleted, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsDeleted), nil
}
//IndexPage represents the content of the index page, available on "/"
//The index page shows a list of all books stored on db
type IndexPage struct {
	AllClients []Client
}

//BookPage represents the content of the book page, available on "/book.html"
//The book page shows info about a given book
type ClientPage struct {
	TargetClient Client
}

type AgainresqPage struct {
	AllCerts []Cert
}
//Book represents a book object
type Client struct {
	ID          int
	Hostname    string
	IP          string
	Dt          string
	Tm          string
	
}

type Cert struct {
	Certname    string
	Selectionval string	
}



//ErrorPage represents shows an error message, available on "/book.html"
type ErrorPage struct {
	ErrorMsg string
}

func handleSaveClient(w http.ResponseWriter, r *http.Request) {
	var id = 0
	var err error

	r.ParseForm()
	params := r.PostForm
	fmt.Println(params)
	idStr := params.Get("id")

	if len(idStr) > 0 {
		id, err = strconv.Atoi(idStr)
		if err != nil {
			renderErrorPage(w, err)
			return
		}
	}

	hostname := params.Get("hostname")
	ip := params.Get("ip")

	pagesStr := params.Get("dt")
	

	tmStr := params.Get("tm")
	
	var tm string

	if len(tmStr) > 0 {
		fmt.Println(tm)
		
	}

	if id == 0 {
		_, err = insertClient(hostname, ip,pagesStr, tm)
	} else {
		_, err = updateClient(id, hostname, ip, pagesStr, tm)
	}

	if err != nil {
		renderErrorPage(w, err)
		return
	}

	http.Redirect(w, r, "/", 302)
}


func handleListClients(w http.ResponseWriter, r *http.Request) {
	clients, err := allClients()
	if err != nil {
		renderErrorPage(w, err)
		return
	}

	buf, err := ioutil.ReadFile("www/index.html")
	if err != nil {
		renderErrorPage(w, err)
		return
	}

	var page = IndexPage{AllClients: clients}
	indexPage := string(buf)
	t := template.Must(template.New("indexPage").Parse(indexPage))
	t.Execute(w, page)
}

func againresq(w http.ResponseWriter, r *http.Request) {
	certs, err := allCerts()
	if err != nil {
		renderErrorPage(w, err)
		return
		//return certs,err
	}

	buf, err := ioutil.ReadFile("www/againresq.html")
	if err != nil {
		renderErrorPage(w, err)
		//return certs,err
		return
	}
	
     /*   w.Header().Set("CONTENT-TYPE", "text/html; charset=UTF-8")          
        
       	rows, err := db.Query("SELECT * FROM certs")
	check(err)
	//var Cert1s string
	//var data []Cert1s
	for rows.Next() {
	        var certname string
	        var selectionval string
		err= rows.Scan(&certname, &selectionval)
		//fmt.Fprintf(w,certname) 
		
		//fmt.Fprintf(w,selectionval)
	

		if err == nil { 
		currentCert := Cert{Certname: certname, Selectionval: selectionval}
				fmt.Println(currentCert)
				//fmt.Fprint(currentCert)
				certs = append(certs, currentCert)
				//return certs,err
				//}
					
			} else {
				//return certs, err
			}
			
			
	} */
        
        //var page3 = AgainresqPage{AllCerts: certs}
      //  json.NewEncoder(w).Encode(page3)
        
	
	againresqPage := string(buf)
	var page = AgainresqPage{AllCerts: certs}
	t := template.Must(template.New("againresqPage").Parse(againresqPage))
	t.Execute(w, page)
}


func handleViewClient(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	idStr := params.Get("id")
	fmt.Println(idStr)

	var currentClient = Client{}
	//currentClient.tm = time.Now()

	if len(idStr) > 0 {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			renderErrorPage(w, err)
		}

		currentClient, err = getClient(id)
		if err != nil {
			renderErrorPage(w, err)
			return
		}
	}

	buf, err := ioutil.ReadFile("www/client.html")
	if err != nil {
		renderErrorPage(w, err)
		return
	}

	var page = ClientPage{TargetClient: currentClient}
	clientPage := string(buf)
	t := template.Must(template.New("clientPage").Parse(clientPage))
	err = t.Execute(w, page)
	if err != nil {
		renderErrorPage(w, err)
		return
	}
}


func handleViewCert(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	idStr := params.Get("id")

	var currentClient = Client{}
	///currentClient.tm = time.Now()

	if len(idStr) > 0 {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			renderErrorPage(w, err)
			return
		}

		currentClient, err = getClient(id)
		if err != nil {
			renderErrorPage(w, err)
			return
		}
	}

	buf, err := ioutil.ReadFile("www/cert.html")
	if err != nil {
		renderErrorPage(w, err)
		return
	}

	var page = ClientPage{TargetClient: currentClient}
	clientPage := string(buf)
	t := template.Must(template.New("clientPage").Parse(clientPage))
	err = t.Execute(w, page)
	if err != nil {
		renderErrorPage(w, err)
		return
	}
}


func check(err error) {
	if err != nil {
		panic(err)
	}
}

type MsgAttachment struct {
	Mattachmentnr sql.NullInt64
	Messagefk     sql.NullInt64
	Aname         sql.NullString
	Mblob         []byte
}

func (a MsgAttachment) Insert(db *sql.DB) error {
	stmt, err := db.Prepare("INSERT INTO testit ( messagefk, aname, ablob) Values($1,$2,$3)")

	if err != nil {
		return fmt.Errorf("cannot prepare: %v", err)
	}
	_, err = stmt.Exec(a.Messagefk, a.Aname, a.Mblob)
	if err != nil {
		return fmt.Errorf("cannot exec: %v", err)
	}
	return nil
}

type MsgAttachment1 struct {
	imgname     sql.NullString
	img         []byte
}

func (b MsgAttachment1) Insert1(db *sql.DB) error {
	stmt, err := db.Prepare("INSERT INTO images (imgname, img) Values($1,$2)")

	if err != nil {
		return fmt.Errorf("cannot prepare: %v", err)
	}
	_, err = stmt.Exec(b.imgname, b.img)
	if err != nil {
		return fmt.Errorf("cannot exec: %v", err)
	}
	return nil
}

        func platcert(w http.ResponseWriter, req *http.Request){

        var currentClient = Client{}
	//currentClient.tm = time.Now()

	buf, err := ioutil.ReadFile("www/platcert.html")
	if err != nil {
		renderErrorPage(w, err)
		return
	}
        var s string
        if req.Method == http.MethodPost {
        f, _, err := req.FormFile("usrfile")
        if err != nil {
                log.Println(err)
                http.Error(w, "Error uploading file", http.StatusInternalServerError)
                return
         }
         defer f.Close()
         
         bs, err := ioutil.ReadAll(f)
         if err != nil{
                  log.Println(err)
                  http.Error(w, "Error reading file", http.StatusInternalServerError)
                  return
                  }
                  s=string(bs)
                  s1:=  MsgAttachment1{
		imgname:         sql.NullString{String: "bname"},
		img:         []byte(s),
	}
                  err = s1.Insert1(db)
	check(err)
                  }
         
                  
        fmt.Println("Sucessfully uploaded a file")          
        w.Header().Set("CONTENT-TYPE", "text/html; charset=UTF-8")          
        
       	rows, err := db.Query("select img from images")
	check(err)
	var data []byte
	for rows.Next() {
		rows.Scan(&data)
			
	} 
	ans:=string(data)
	fmt.Println(ans)
	fmt.Fprintf(w, `
        <h4></h4>`+ans)
        var page = ClientPage{TargetClient: currentClient}
	clientPage := string(buf)
	t := template.Must(template.New("clientPage").Parse(clientPage))
	err = t.Execute(w, page)
	if err != nil {
		renderErrorPage(w, err)
		return
	}
        
        }  
 
        
        
        func endocert(w http.ResponseWriter, req *http.Request){

        var currentClient = Client{}
	//currentClient.tm = time.Now()

	buf, err := ioutil.ReadFile("www/endocert.html")
	if err != nil {
		renderErrorPage(w, err)
		return
	}
        var s string
        if req.Method == http.MethodPost {
        f, _, err := req.FormFile("usrfile")
        if err != nil {
                log.Println(err)
                http.Error(w, "Error uploading file", http.StatusInternalServerError)
                return
         }
         defer f.Close()
         
         bs, err := ioutil.ReadAll(f)
         if err != nil{
                  log.Println(err)
                  http.Error(w, "Error reading file", http.StatusInternalServerError)
                  return
                  }
                  s=string(bs)
                  s1:=  MsgAttachment1{
		imgname:         sql.NullString{String: "bname"},
		img:         []byte(s),
	}
                  err = s1.Insert1(db)
	check(err)
                  }
         
                  
        fmt.Println("Sucessfully uploaded a file")          
        w.Header().Set("CONTENT-TYPE", "text/html; charset=UTF-8")          
        
       	rows, err := db.Query("select img from images")
	check(err)
	var data []byte
	for rows.Next() {
		rows.Scan(&data)
			
	} 
	ans:=string(data)
	fmt.Println(ans)
	fmt.Fprintf(w, `
        <h4></h4>`+ans)
        var page = ClientPage{TargetClient: currentClient}
	clientPage := string(buf)
	t := template.Must(template.New("clientPage").Parse(clientPage))
	err = t.Execute(w, page)
	if err != nil {
		renderErrorPage(w, err)
		return
	}
        
        }  
        
        
        func eventlog(w http.ResponseWriter, req *http.Request){

        var currentClient = Client{}
	//currentClient.tm = time.Now()

	buf, err := ioutil.ReadFile("www/eventlog.html")
	if err != nil {
		renderErrorPage(w, err)
		return
	}
        var s string
        if req.Method == http.MethodPost {
        f, _, err := req.FormFile("usrfile")
        if err != nil {
                log.Println(err)
                http.Error(w, "Error uploading file", http.StatusInternalServerError)
                return
         }
         defer f.Close()
         
         bs, err := ioutil.ReadAll(f)
         if err != nil{
                  log.Println(err)
                  http.Error(w, "Error reading file", http.StatusInternalServerError)
                  return
                  }
                  s=string(bs)
                  s1:=  MsgAttachment1{
		imgname:         sql.NullString{String: "bname"},
		img:         []byte(s),
	}
                  err = s1.Insert1(db)
	check(err)
                  }
         
                  
        fmt.Println("Sucessfully uploaded a file")          
        w.Header().Set("CONTENT-TYPE", "text/html; charset=UTF-8")          
        
       	rows, err := db.Query("select img from images")
	check(err)
	var data []byte
	for rows.Next() {
		rows.Scan(&data)
			
	} 
	ans:=string(data)
	fmt.Println(ans)
	fmt.Fprintf(w, `
        <h4></h4>`+ans)
        var page = ClientPage{TargetClient: currentClient}
	clientPage := string(buf)
	t := template.Must(template.New("clientPage").Parse(clientPage))
	err = t.Execute(w, page)
	if err != nil {
		renderErrorPage(w, err)
		return
	}
        
        }      
        
        func imalogs(w http.ResponseWriter, req *http.Request){

        var currentClient = Client{}
	//currentClient.tm = time.Now()

	buf, err := ioutil.ReadFile("www/imalogs.html")
	if err != nil {
		renderErrorPage(w, err)
		return
	}
        var s string
        if req.Method == http.MethodPost {
        f, _, err := req.FormFile("usrfile")
        if err != nil {
                log.Println(err)
                http.Error(w, "Error uploading file", http.StatusInternalServerError)
                return
         }
         defer f.Close()
         
         bs, err := ioutil.ReadAll(f)
         if err != nil{
                  log.Println(err)
                  http.Error(w, "Error reading file", http.StatusInternalServerError)
                  return
                  }
                  s=string(bs)
                  s1:=  MsgAttachment1{
		imgname:         sql.NullString{String: "bname"},
		img:         []byte(s),
	}
                  err = s1.Insert1(db)
	check(err)
                  }
         
                  
        fmt.Println("Sucessfully uploaded a file")          
        w.Header().Set("CONTENT-TYPE", "text/html; charset=UTF-8")          
        
       	rows, err := db.Query("select img from images")
	check(err)
	var data []byte
	for rows.Next() {
		rows.Scan(&data)
			
	} 
	ans:=string(data)
	fmt.Println(ans)
	fmt.Fprintf(w, `
        <h4></h4>`+ans)
        var page = ClientPage{TargetClient: currentClient}
	clientPage := string(buf)
	t := template.Must(template.New("clientPage").Parse(clientPage))
	err = t.Execute(w, page)
	if err != nil {
		renderErrorPage(w, err)
		return
	}
        
        }  

// ProcessCheckboxes will process checkboxes
func ProcessCheckboxes(w http.ResponseWriter, r *http.Request) {
    var currentClient = Client{}
	//currentClient.tm = time.Now()

	buf, err := ioutil.ReadFile("www/aga.html")
	//fmt.Println("Sucessfully uploaded a file")    
	if err != nil {
		renderErrorPage(w, err)
		return
	}
	
    r.ParseForm()
  // fmt.Printf("%+v\n", r.Form)

    productsSelected := r.Form["ns_license"]
    fmt.Println(productsSelected)
    //log.Println(contains(productsSelected, "NeuroSolutions"))
    
    var page = ClientPage{TargetClient: currentClient}
	clientPage := string(buf)
	t := template.Must(template.New("clientPage").Parse(clientPage))
	err = t.Execute(w, page)
	if err != nil {
		renderErrorPage(w, err)
		return
	}
	//fmt.Println(productsSelected)
}

func contains(slice []string, item string) bool {
  
 set := make(map[string]struct{}, len(slice))
    for _, s := range slice {
        set[s] = struct{}{}
    }
    _, ok := set[item]
    return ok
}
        
func again(w http.ResponseWriter, r *http.Request) {
        var platcert string
        platcert = "plat1"
       // fmt.Println("Avani", platcert)
        //fmt.Println(platcert)
	var currentClient = Client{}
	//currentClient.tm = time.Now()

	buf, err := ioutil.ReadFile("www/aga.html")
	//fmt.Println("Sucessfully uploaded a file")    
	if err != nil {
		renderErrorPage(w, err)
		return
	}
	
	r.ParseForm()
	//temp:= r.PostForm
	productsSelected := r.Form.Get("ns_license")
	//productsSelected := r.Form.Get("selection")
       fmt.Println("test:", productsSelected)
        //t1:= temp.Get("ns_license")
       // fmt.Fprintln(w,productsSelected)
	//t1:= temp.Get("ns_license")
	//fmt.Fprintln(w, t1)
	//inserting(certname, productsSelected) 
	
	//db.QueryRow(`INSERT INTO platcerts(selectionval) VALUES($1)`, productsSelected)
	db.QueryRow(`UPDATE cert1 set certname=$1, selectionval=$2 where id=3`,platcert, productsSelected )
	//db.QueryRow(`INSERT INTO cert1(selectionval) VALUES($1)`, productsSelected)
	//db.QueryRow(`INSERT selectionval INTO certs where certname ='$1'` ,platcert).Scan(platcert, productsSelected)
	
//id, hostname, ip, dt, tm FROM clients where id = $1`, clientID).Scan(&id, &hostname, &ip, &dt, &tm)
         
//UPDATE  cert1 set certname='t1', selectionval='1' where id=1;


        //rows,err := db.Query(`SELECT selectionval FROM certs`);
rows,err := db.Query(`SELECT selectionval from cert1 where id=3`);
if err != nil {
    log.Fatal(err)
}     


for rows.Next() {
    var value string
    if err := rows.Scan(&value); err != nil {
        log.Fatal(err)
    }
   // fmt.Printf("Value: %t\n",value);
    fmt.Println(value);
}

      //  certview()
	var page = ClientPage{TargetClient: currentClient}
	clientPage := string(buf)
	t := template.Must(template.New("clientPage").Parse(clientPage))
	err = t.Execute(w, page)
	if err != nil {
		renderErrorPage(w, err)
		return
	}
	
	
}

func againres(w http.ResponseWriter, r *http.Request) {
        var platcert string
        platcert = "plat1"
       // fmt.Println("Avani", platcert)
        //fmt.Println(platcert)
	var currentClient = Client{}
	//currentClient.tm = time.Now()

	buf, err := ioutil.ReadFile("www/aga.html")
	//fmt.Println("Sucessfully uploaded a file")    
	if err != nil {
		renderErrorPage(w, err)
		return
	}
	
	r.ParseForm()
	//temp:= r.PostForm
	productsSelected := r.Form.Get("ns_license")
	//productsSelected := r.Form.Get("selection")
       fmt.Println("test:", productsSelected)
        //t1:= temp.Get("ns_license")
       // fmt.Fprintln(w,productsSelected)
	//t1:= temp.Get("ns_license")
	//fmt.Fprintln(w, t1)
	//inserting(certname, productsSelected) 
	
	//db.QueryRow(`INSERT INTO platcerts(selectionval) VALUES($1)`, productsSelected)
	db.QueryRow(`UPDATE cert1 set certname=$1, selectionval=$2 where id=3`,platcert, productsSelected )
	//db.QueryRow(`INSERT INTO cert1(selectionval) VALUES($1)`, productsSelected)
	//db.QueryRow(`INSERT selectionval INTO certs where certname ='$1'` ,platcert).Scan(platcert, productsSelected)
	
//id, hostname, ip, dt, tm FROM clients where id = $1`, clientID).Scan(&id, &hostname, &ip, &dt, &tm)
         
//UPDATE  cert1 set certname='t1', selectionval='1' where id=1;


        //rows,err := db.Query(`SELECT selectionval FROM certs`);
rows,err := db.Query(`SELECT selectionval from cert1 where id=3`);
if err != nil {
    log.Fatal(err)
}     


for rows.Next() {
    var value string
    if err := rows.Scan(&value); err != nil {
        log.Fatal(err)
    }
   // fmt.Printf("Value: %t\n",value);
    fmt.Println(value);
}

      //  certview()
	var page = ClientPage{TargetClient: currentClient}
	clientPage := string(buf)
	t := template.Must(template.New("clientPage").Parse(clientPage))
	err = t.Execute(w, page)
	if err != nil {
		renderErrorPage(w, err)
		return
	}
	
	
}
                    
func handleViewPolicy(w http.ResponseWriter, r *http.Request) {
      
	var currentClient = Client{}
	//currentClient.tm = time.Now()

	buf, err := ioutil.ReadFile("www/policy.html")
	//fmt.Println("Sucessfully uploaded a file")    
	if err != nil {
		renderErrorPage(w, err)
		return
	}
	
	r.ParseForm()
	temp:= r.PostForm
	//fmt.Fprintln(w, temp)
	t1:= temp.Get("selection")
	fmt.Fprintln(w, t1)
	//var platcert string 
	//var s1 string
	//slice:=[]string{`true`,`false`}
       // if r.Method == http.MethodPost {
       // fmt.Println("Sucessfully uploaded a file")  
       // r.ParseForm()
	//params := r.PostForm
	//fmt.Println(params)
       // selection := params.Get("selection1")
      // for _, v := range slice {
       //var t1 string
      // var t1 = r.Form.Get("selection1") 
     //  fmt.Println(t1)
   //   if v == r.Form.Get("selection1") {
       // fmt.Println(r.FormValue("selection1")) 
       //fmt.Println(r.url.Values("myModal1")) 
      // return true
     //  }
     //  }
      // fmt.Println() 
       //fmt.Println(selection)
        //  fmt.Println("Sucessfully uploaded a file1")     
      //  bs1, err := ioutil.ReadAll(selection1)
      //   if err != nil{
            //      log.Println(err)
             //     http.Error(w, "Error reading file", http.StatusInternalServerError)
            //      return
              //    }
             //     s1=string(bs1)
        //  db.Exec(`INSERT INTO certs(certname, selectionval) VALUES($1, $2)`, platcert, v)
               //  err= insertcertselection(platcert1, s1)
	//check(err)
	
	//}
	
	/*m := MsgAttachment{
		Mattachmentnr: sql.NullInt64{Int64: 12},
		Messagefk:     sql.NullInt64{Int64: 99},
		Aname:         sql.NullString{String: "a name"},
		Mblob:         []byte("hello \xff\x00world"),
	}
	err = m.Insert(db)
	check(err)

	m = MsgAttachment{
		Mattachmentnr: sql.NullInt64{Int64: 12},
		Messagefk:     sql.NullInt64{Int64: 99},
		Aname:         sql.NullString{String: "a name"},
		Mblob:         []byte("more data \x85\xd5"),
	}
	err = m.Insert(db)
	m = MsgAttachment{
		Mattachmentnr: sql.NullInt64{Int64: 12},
		Messagefk:     sql.NullInt64{Int64: 99},
		Aname:         sql.NullString{String: "a name"},
		Mblob:         []byte("more data \x85\xd5"),
	}
	err = m.Insert(db)
	check(err)
	rows, err := db.Query("select ablob from testit")
	check(err)
	var data []byte
	for rows.Next() {
		rows.Scan(&data)
		fmt.Printf("%q\n", data)
	}*/
	var page = ClientPage{TargetClient: currentClient}
	clientPage := string(buf)
	t := template.Must(template.New("clientPage").Parse(clientPage))
	err = t.Execute(w, page)
	if err != nil {
		renderErrorPage(w, err)
		return
	}
}


func handleViewReport(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	idStr := params.Get("id")

	var currentClient = Client{}
	//currentClient.tm = time.Now()

	if len(idStr) > 0 {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			renderErrorPage(w, err)
			return
		}

		currentClient, err = getClient(id)
		if err != nil {
			renderErrorPage(w, err)
			return
		}
	}

	buf, err := ioutil.ReadFile("www/report.html")
	if err != nil {
		renderErrorPage(w, err)
		return
	}

	var page = ClientPage{TargetClient: currentClient}
	clientPage := string(buf)
	t := template.Must(template.New("clientPage").Parse(clientPage))
	err = t.Execute(w, page)
	if err != nil {
		renderErrorPage(w, err)
		return
	}
}

func handleDeleteClient(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	idStr := params.Get("id")

	if len(idStr) > 0 {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			renderErrorPage(w, err)
			return
		}

		n, err := removeClient(id)
		if err != nil {
			renderErrorPage(w, err)
			return
		}

		fmt.Printf("Rows removed: %v\n", n)
	}
	http.Redirect(w, r, "/", 302)
}

func renderErrorPage(w http.ResponseWriter, errorMsg error) {
	buf, err := ioutil.ReadFile("www/error.html")
	if err != nil {
		log.Printf("%v\n", err)
		fmt.Fprintf(w, "%v\n", err)
		return
	}

	var page = ErrorPage{ErrorMsg: errorMsg.Error()}
	errorPage := string(buf)
	t := template.Must(template.New("errorPage").Parse(errorPage))
	t.Execute(w, page)
}

func init() {
	tmpDB, err := sql.Open("postgres", "user=postgres password=myPassword dbname=books_database sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	db = tmpDB
	fmt.Println("Sucessfully connected to DB!!!")
}



func main() {
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("www/assets"))))

	http.HandleFunc("/", handleListClients)
	http.HandleFunc("/client.html", handleViewClient)
	http.HandleFunc("/cert.html", handleViewCert)
	http.HandleFunc("/policy.html", handleViewPolicy)
	http.HandleFunc("/aga.html", again)
	http.HandleFunc("/againresq.html",againresq)
	http.HandleFunc("/platcert.html", platcert)
	http.HandleFunc("/endocert.html", endocert)
	http.HandleFunc("/eventlog.html", eventlog)
	http.HandleFunc("/imalogs.html", imalogs)
	http.HandleFunc("/report.html", handleViewReport)
       http.HandleFunc("/save",handleSaveClient)
	http.HandleFunc("/delete", handleDeleteClient)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
