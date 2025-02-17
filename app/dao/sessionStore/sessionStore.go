package sessionStore

import (
	"context"
	"fmt"
	"log"
	"strings"

	audit "github.com/mt1976/frantic-core/dao/audit"
	"github.com/mt1976/frantic-core/dao/database"
	io "github.com/mt1976/frantic-core/ioHelpers"
	logger "github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/paths"
	stopwatch "github.com/mt1976/frantic-core/timing"
)

func (s *Aegis_SessionStore) Update(ctx context.Context, note string) error {
	clock := stopwatch.Start(domain, "Update", s.ID)
	if err := s.validate(); err != nil {
		clock.Stop(0)
		return err
	}

	if calculationError := s.calculate(); calculationError != nil {
		logger.ErrorLogger.Printf("[%v] Calculating %v", strings.ToUpper(domain), calculationError)
		clock.Stop(0)
		return calculationError
	}

	if _, validationError := s.prepare(); validationError != nil {
		logger.ErrorLogger.Printf("[%v] Validating %v", strings.ToUpper(domain), validationError.Error())
		clock.Stop(0)
		return validationError
	}

	_ = s.Audit.Action(ctx, audit.UPDATE.WithMessage(note))

	if err := database.Update(s); err != nil {
		msg := fmt.Sprintf("[%v] Updating [%v]", strings.ToUpper(domain), err.Error())
		log.Panic(msg)
	}

	logger.AuditLogger.Printf("[%v] [%v] ID=[%v] Notes[%v]", audit.UPDATE, strings.ToUpper(domain), s.ID, note)
	clock.Stop(1)

	return nil
}

func GetById(id string) (Aegis_SessionStore, error) {
	return GetBy("ID", id)
}

func GetBy(field, value string) (Aegis_SessionStore, error) {
	clock := stopwatch.Start(domain, "Get", value)
	u := Aegis_SessionStore{}
	if err := database.Retrieve(field, value, &u); err != nil {
		clock.Stop(0)
		return Aegis_SessionStore{}, fmt.Errorf("[%v] Reading Id=[%v][%v]", strings.ToUpper(domain), value, err.Error())
	}
	if err := u.PostGet(); err != nil {
		clock.Stop(0)
		return Aegis_SessionStore{}, err
	}
	clock.Stop(1)
	return u, nil
}

func GetAll() ([]Aegis_SessionStore, error) {
	uList := []Aegis_SessionStore{}
	clock := stopwatch.Start(domain, "Get All", "")
	if errG := database.GetAll(&uList); errG != nil {
		logger.ErrorLogger.Printf("[%v] Reading Id=[%v] %v ", strings.ToUpper(domain), "ALL", errG.Error())
		panic(errG)
	}

	dList, errPost := PostGet(&uList)
	if errPost != nil {
		clock.Stop(0)
		return nil, errPost
	}

	clock.Stop(len(uList))
	return dList, nil
}

func (s *Aegis_SessionStore) Delete(ctx context.Context, note string) error {
	clock := stopwatch.Start(domain, "Delete", s.ID)
	logger.AuditLogger.Printf("DEL: [%v] Id=[%v]", strings.ToUpper(domain), s.ID)
	_ = s.Audit.Action(ctx, audit.DELETE.WithMessage(note))
	s.Dump("DEL")

	if err := database.Delete(s); err != nil {
		logger.ErrorLogger.Printf("[%v] Deleting %v ", strings.ToUpper(domain), err)
		panic(err)
	}

	logger.AuditLogger.Printf("DEL: [%v] ID=[%v] ", strings.ToUpper(domain), s.ID)
	clock.Stop(1)
	return nil
}

func DeleteByID(ctx context.Context, id, note string) error {
	clock := stopwatch.Start(domain, "Delete", id)
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
	clock.Stop(1)
	return nil
}

func (s *Aegis_SessionStore) Spew() {
	ID := fmt.Sprintf("%04v", s.ID)
	if s.ID == "" {
		ID = "NEW"
	}
	logger.InfoLogger.Printf(" [%v] ID=[%v]", strings.ToUpper(domain), ID)
}

func (s *Aegis_SessionStore) validate() error {
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

func (s *Aegis_SessionStore) PostGet() error {
	return nil
}

func (s *Aegis_SessionStore) Dump(name string) {
	io.Dump(domain, paths.Dumps(), name, s.ID, s)
}

func Export() {
	clock := stopwatch.Start(domain, "Export", "")
	dList, _ := GetAll()
	if len(dList) == 0 {
		logger.EventLogger.Printf("[%v] Backup [%v] no data found", strings.ToUpper(domain), domain)
		clock.Stop(0)
		return
	}
	for _, yy := range dList {
		yy.Dump("EXPORT")
	}
	clock.Stop(len(dList))
}
