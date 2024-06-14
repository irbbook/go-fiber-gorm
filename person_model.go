package main

import (
	"gorm.io/gorm"
)

type Person struct {
	gorm.Model
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Gender    string `json:"gender"`
}

func createPerson(db *gorm.DB, person *Person) error {
	result := db.Create(person)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func getPeople(db *gorm.DB) ([]Person, error) {
	var people []Person

	result := db.Find(&people)

	if result.Error != nil {
		return []Person{}, result.Error
	}

	return people, nil
}

func getPerson(db *gorm.DB, id uint) (*Person, error) {
	var person Person

	result := db.First(&person, id)

	if result.Error != nil {
		return &Person{}, result.Error
	}

	return &person, nil
}

func updatePerson(db *gorm.DB, person *Person) (*Person, error) {
	currentPerson := new(Person)
	result := db.Where("id = ?", person.ID).First(currentPerson)
	if result.Error != nil {
		return &Person{}, result.Error
	}

	currentPerson.Firstname = person.Firstname
	currentPerson.Lastname = person.Lastname
	currentPerson.Gender = person.Gender

	result = db.Save(currentPerson)

	if result.Error != nil {
		return &Person{}, result.Error
	}

	return person, nil
}

func deletePerson(db *gorm.DB, id uint) error {
	result := db.Where("id = ?", id).First(&Person{})
	if result.Error != nil {
		return result.Error
	}

	result = db.Delete(&Person{}, id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
