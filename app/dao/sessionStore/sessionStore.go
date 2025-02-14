package sessionStore

import (
	"context"
	"fmt"
	"log"
	"strings"

	audit "github.com/mt1976/frantic-core/dao/audit"
	"github.com/mt1976/frantic-core/dao/database"
	"github.com/mt1976/frantic-core/io"
	"github.com/mt1976/frantic-core/logger"
	"github.com/mt1976/frantic-core/paths"
	stopwatch "github.com/mt1976/frantic-core/timing"
)

var sessionExpiry = 30 // minutes

func (u *Aegis_SessionStore) Update(ctx context.Context, note string) error {
	if err := u.validate(); err != nil {
		return err
	}

	if calculationError := u.calculate(); calculationError != nil {
		logger.ErrorLogger.Printf("[%v] Calculating %v", strings.ToUpper(domain), calculationError)
		return calculationError
	}

	if _, validationError := u.prepare(); validationError != nil {
		logger.ErrorLogger.Printf("[%v] Validating %v", strings.ToUpper(domain), validationError.Error())
		return validationError
	}

	_ = u.Audit.Action(ctx, audit.UPDATE.WithMessage(note))

	if err := database.Update(u); err != nil {
		msg := fmt.Sprintf("[%v] Updating [%v]", strings.ToUpper(domain), err.Error())
		log.Panic(msg)
	}

	logger.AuditLogger.Printf("[%v] [%v] ID=[%v] Notes[%v]", audit.UPDATE, strings.ToUpper(domain), u.ID, note)

	return nil
}

func GetById(id string) (Aegis_SessionStore, error) {
	return GetBy("ID", id)
}

func GetBy(field, value string) (Aegis_SessionStore, error) {
	get := stopwatch.Start(domain, "Get", value)
	u := Aegis_SessionStore{}
	if err := database.Retrieve(field, value, &u); err != nil {
		return Aegis_SessionStore{}, fmt.Errorf("[%v] Reading Id=[%v][%v]", strings.ToUpper(domain), value, err.Error())
	}
	if err := u.PostGet(); err != nil {
		return Aegis_SessionStore{}, err
	}
	get.Stop(1)
	return u, nil
}

func GetAll() ([]Aegis_SessionStore, error) {
	uList := []Aegis_SessionStore{}
	gall := stopwatch.Start(domain, "Get All", "Get")
	if errG := database.GetAll(&uList); errG != nil {
		logger.ErrorLogger.Printf("[%v] Reading Id=[%v] %v ", strings.ToUpper(domain), "ALL", errG.Error())
		panic(errG)
	}

	dList, errPost := PostGet(&uList)
	if errPost != nil {
		return nil, errPost
	}

	gall.Stop(len(uList))
	return dList, nil
}

func (u *Aegis_SessionStore) Delete(ctx context.Context, note string) error {
	logger.AuditLogger.Printf("DEL: [%v] Id=[%v]", strings.ToUpper(domain), u.ID)
	_ = u.Audit.Action(ctx, audit.DELETE.WithMessage(note))
	u.Dump("DEL")

	if err := database.Drop(u); err != nil {
		logger.ErrorLogger.Printf("[%v] Deleting %v ", strings.ToUpper(domain), err)
		panic(err)
	}

	logger.AuditLogger.Printf("DEL: [%v] ID=[%v] ", strings.ToUpper(domain), u.ID)
	return nil
}

func DeleteByID(ctx context.Context, id, note string) error {
	logger.AuditLogger.Printf("DLI: [%v] Id=[%v]", domain, id)
	dest, err := GetById(id)
	if err != nil {
		logger.ErrorLogger.Printf("[%v] Reading Id=[%v] %v", strings.ToUpper(domain), id, err.Error())
		panic(err)
	}

	_ = dest.Audit.Action(ctx, audit.DELETE.WithMessage(note))
	dest.Dump("DEL")

	if err := database.Delete(&dest); err != nil {
		logger.ErrorLogger.Printf("[%v] Deleting %v ", strings.ToUpper(domain), err)
		panic(err)
	}

	logger.AuditLogger.Printf("DLI: [%v] Id=[%v] ", strings.ToUpper(domain), id)
	return nil
}

func (u *Aegis_SessionStore) Spew() {
	ID := fmt.Sprintf("%04v", u.ID)
	if u.ID == "" {
		ID = "NEW"
	}
	logger.InfoLogger.Printf(" [%v] ID=[%v]", strings.ToUpper(domain), ID)
}

func (u *Aegis_SessionStore) validate() error {
	return nil
}

func PostGet(userList *[]Aegis_SessionStore) ([]Aegis_SessionStore, error) {
	newList := []Aegis_SessionStore{}
	for _, user := range *userList {
		if err := user.PostGet(); err != nil {
			return nil, err
		}
		newList = append(newList, user)
	}
	return newList, nil
}

func (u *Aegis_SessionStore) PostGet() error {
	return nil
}

func (u *Aegis_SessionStore) Dump(name string) {
	io.Dump(domain, paths.Dumps(), name, u.ID, u)
}

func Export() {
	dList, _ := GetAll()
	if len(dList) == 0 {
		logger.EventLogger.Printf("[%v] Backup [%v] no data found", strings.ToUpper(domain), domain)
		return
	}
	for _, yy := range dList {
		yy.Dump("EXPORT")
	}
}
