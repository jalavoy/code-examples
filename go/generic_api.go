// jalavoy - 06.08.2018
// this is a generic data collection and retrieval API to a MySQL database
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// webserver settings
const listenAddr string = "127.0.0.1"
const listenPort string = "8080"

// sql settings
const sqlUser string = "api"
const sqlPass string = "password"
const sqlDB string = "api"

type submit struct {
	Hostdomain string
	Ips        map[string][]string
	Hostname   string
	Primaryip  string
}

func main() {
	// setup our router
	router := mux.NewRouter()

	// setup our listener
	srv := &http.Server{
		Addr:         listenAddr + ":" + listenPort,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	// handle routes
	router.HandleFunc("/", defaultHandler)
	router.HandleFunc("/get", getHandler).Methods("GET")
	router.HandleFunc("/submit", submitHandler).Methods("POST")

	// start the listener as a child so we don't cause blocking
	log.Printf("Starting server on %s:%s", listenAddr, listenPort)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// capture SIGINT so we can shutdown gracefully
	signal.Notify(c, os.Interrupt)

	// block until we get a signal
	<-c

	// set a timeout
	context, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	srv.Shutdown(context)
	log.Println("Shutting down")
	os.Exit(0)

}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := `SELECT 
			servers.hostname, 
			INET_NTOA(servers.primary_ip), 
			domains.domain_name, 
			INET_NTOA(ip_addresses.ip_address)
		FROM servers 
		INNER JOIN domains 
			ON servers.id = domains.server_id
		INNER JOIN ip_addresses
			ON servers.id = ip_addresses.server_id`
	if r.FormValue("domain") != "" {
		query = fmt.Sprintf("%s WHERE domains.domain_name = '%s'", query, r.FormValue("domain"))
	} else if r.FormValue("ip") != "" {
		query = fmt.Sprintf("%s WHERE servers.primary_ip = INET_ATON(%s)", query, r.FormValue("ip"))
	} else if r.FormValue("server") != "" {
		query = fmt.Sprintf("%s WHERE servers.hostname = '%s'", query, r.FormValue("server"))
	} else {
		// else no input specified
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db, err := connectDB()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	var (
		hostname   string
		primaryIP  string
		domainName string
		siteIP     string
	)
	err = db.QueryRow(query).Scan(&hostname, &primaryIP, &domainName, &siteIP)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	var results = map[string]string{
		"hostname":   hostname,
		"primaryIP":  primaryIP,
		"domainName": domainName,
		"siteIP":     siteIP,
	}
	json, err := json.Marshal(results)
	if err != nil {
		log.Panic(err)
		return
	}
	w.Write(json)
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	var s submit
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	err = json.Unmarshal(b, &s)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		err := insertDB(s)
		if err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusAccepted)
	}
}

func connectDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", sqlUser+":"+sqlPass+"@/"+sqlDB)
	if err != nil {
		log.Panic(err)
		return nil, err
	}
	return db, nil
}

func insertDB(s submit) error {
	db, err := connectDB()
	defer db.Close()
	if err != nil {
		log.Panic(err)
		return err
	}

	// servers table
	sth, err := db.Prepare("INSERT INTO servers ( hostname, primary_ip ) VALUES ( ?, INET_ATON(?) ) ON DUPLICATE KEY UPDATE primary_ip = INET_ATON(?)")
	if err != nil {
		log.Panic(err)
		return err
	}
	defer sth.Close()

	res, err := sth.Exec(s.Hostname, s.Primaryip, s.Primaryip)
	if err != nil {
		log.Panic(err)
		return err
	}
	serverID, err := res.LastInsertId()
	if err != nil {
		log.Panic(err)
		return err
	}
	if serverID == 0 {
		err = db.QueryRow("SELECT id FROM servers WHERE hostname = ? AND primary_ip = INET_ATON(?)", s.Hostname, s.Primaryip).Scan(&serverID)
		if err != nil || serverID == 0 {
			log.Panic(err)
			return err
		}
	}

	for ip, domainArray := range s.Ips {
		// ip_addresses table
		sth, err = db.Prepare("INSERT INTO ip_addresses ( ip_address, server_id ) VALUES (INET_ATON(?), ? ) ON DUPLICATE KEY UPDATE server_id = ?")
		if err != nil {
			log.Panic(err)
			return err
		}

		res, err = sth.Exec(ip, serverID, serverID)
		ipID, err := res.LastInsertId()
		if err != nil {
			log.Panic(err)
			return err
		}
		if ipID == 0 {
			err = db.QueryRow("SELECT id FROM ip_addresses WHERE ip_address = INET_ATON(?) AND server_id = ?", ip, serverID).Scan(&ipID)
			if err != nil || ipID == 0 {
				log.Panic(err)
				return err
			}
		}
		// domain table
		for _, domain := range domainArray {
			sth, err = db.Prepare("INSERT INTO domains ( server_id, ip_id, domain_name ) VALUES ( ?, ?, ? ) ON DUPLICATE KEY UPDATE server_id = ?, ip_id = ?")
			if err != nil {
				log.Panic(err)
				return err
			}

			_, err = sth.Exec(serverID, ipID, domain, serverID, ipID)
			if err != nil {
				log.Panic(err)
				return err
			}
		}
	}

	return nil
}
