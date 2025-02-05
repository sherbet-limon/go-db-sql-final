package main

import (
	"database/sql"
	"fmt"
	"errors"
	
	
	_ "modernc.org/sqlite"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
res, err := s.db.Exec("INSERT INTO parcel (Client, Status, Address, Created_at) VALUES (:Client, :Status, :Address, :Created_at)",
	sql.Named("Client", p.Client),
	sql.Named("Status", p.Status),
	sql.Named("Address", p.Address),
	sql.Named("Created_at", p.Created_at))
	if err != nil {
        fmt.Println(err)
        return 0, err
	}
	
	id, err := res.LastInsertId()
	if err != nil {
		return 0, nil
	}
	// верните идентификатор последней добавленной записи
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	rows, err := s.db.Query("SELECT Number, Client, Status, Address, Created_at FROM parcel WHERE number = :number", 
	sql.Named("number", number))
    if err != nil {
        return Parcel{}, err
    }
	defer rows.Close()
	// заполните объект Parcel данными из таблицы
	var p Parcel
	if rows.Next() {
        err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.Created_at)
        if err != nil {
            return Parcel{}, err
        }
		return p, nil
    }
		return Parcel{}, err
} 
	


func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	rows, err := s.db.Query("SELECT Number, Client, Status, Address, Created_at FROM parcel WHERE client = :client", 
	sql.Named("client", client))
    if err != nil {
        return nil, err
    }
	defer rows.Close()

	// заполните срез Parcel данными из таблицы
	var res []Parcel
	for rows.Next() {
		var p Parcel
	err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.Created_at)
	if err != nil {
		return nil, err
	}
	res = append(res, p)
}
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	res, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number", 
	sql.Named("status", status),
	sql.Named("number", number))
    if err != nil {
        return err
    }
	fmt.Println(res.LastInsertId())
    fmt.Println(res.RowsAffected())
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	rows, err := s.db.Query("SELECT Status FROM parcel WHERE number = :number", 
		sql.Named("number", number))
    if err != nil {
        return err
    }
	defer rows.Close()

	var actualStatus string
	if rows.Next() {
		if err := rows.Scan(&actualStatus); err != nil {
        return err
		}
	} else {
		return errors.New("посылка не найдена")
	}

	if actualStatus != ParcelStatusRegistered {
		return errors.New("нельзя изменить адрес")
	}
	_, err = s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number", 
		sql.Named("address", address),
		sql.Named("number", number))
	
    return err
}



func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered

	rows, err := s.db.Query("SELECT Status FROM parcel WHERE number = :number", 
		sql.Named("number", number))
    if err != nil {
        return err
    }
	defer rows.Close()

	var actualStatus string
	if rows.Next() {
		if err := rows.Scan(&actualStatus); err != nil {
			return err
		}
	} else {
		return errors.New("посылка не найдена")
	}

	if actualStatus != ParcelStatusRegistered {
		return errors.New("нельзя удалить посылку")
	}
	_, err = s.db.Exec("DELETE FROM parcel WHERE number = :number",
		sql.Named("number", number))

	return err
}