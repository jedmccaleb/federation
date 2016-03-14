package federation

import (
	"fmt"
	"log"
	"regexp"

	"github.com/stellar/federation/db"
	"github.com/stellar/go-stellar-base/keypair"
)


/*
line 71 of request_handler.go

if err != nil {
		if err.Error() == "sql: no rows in result set" {
			////////////////////////// for tip bot ////////////////////
			created := rh.createAccount(name, &record)
			if created == false {
				log.Print("Federation record NOT found")
				http.Error(w, ErrorResponseString("not_found", "Account not found"), http.StatusNotFound)
				return
			} else {
				log.Print("Federation record created")
			}
			//////////////////////////////////////////////////////////
		} else {
			log.Print("Server error: ", err)
			log.Print("Q:" + rh.config.FederationQuery + " " + name)
			http.Error(w, ErrorResponseString("server_error", "Server error"), http.StatusInternalServerError)
			return
		}
	}
 */

////////////////////////// for tip bot ////////////////////
// TODO: Check if github name exists before creating a federated name for it
func (rh *RequestHandler) createAccount(name string, record *db.FederationRecord) bool {

	log.Println("Creating a new account: " + name)
	pattern := "^[a-zA-Z0-9](?:-?[a-zA-Z0-9]){0,38}$"
	exp := regexp.MustCompile(pattern)

	if exp.MatchString(name) == false {
		log.Println("improper name")
		return false
	}

	key, err := keypair.Random()

	record.AccountId = key.Address()
	secretKey := key.Seed()

	sql := fmt.Sprintf("INSERT INTO TipUser (GithubName,AccountID,SecretKey) values ('%s','%s','%s')", name, record.AccountId, secretKey)

	//log.Println(sql)

	_, err = rh.driver.Exec(sql)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
