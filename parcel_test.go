package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"
	"fmt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		Created_at: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	// prepare
	// настройте подключение к БД
	db, err := sql.Open("sqlite", "tracker.db")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()
	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора

	var id int
	id, err = store.Add(parcel)
    require.NoError(t, err)
	require.NotEmpty(t, id)

	// get
	// получите только что добавленную посылку, убедитесь в отсутствии ошибки
	var returnParcel Parcel
    returnParcel, err = store.Get(id)
	require.NoError(t, err)
	require.NotEmpty(t, returnParcel)
    // Проверяем, что данные посылки совпадают с ожидаемыми
	// проверьте, что значения всех полей в полученном объекте совпадают со значениями полей в переменной parcel
    assert.Equal(t, parcel.Client, returnParcel.Client, "Не совпадает")
    assert.Equal(t, parcel.Status, returnParcel.Status, "Не совпадает")
    assert.Equal(t, parcel.Address, returnParcel.Address, "Не совпадает")
	pTime, _ := time.Parse(time.RFC3339, parcel.Created_at)
	rTime, _ := time.Parse(time.RFC3339, returnParcel.Created_at)
	assert.WithinDuration(t, pTime, rTime, time.Second, "Не совпадает")


	// delete
	// удалите добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что посылку больше нельзя получить из БД
	err = store.Delete(id)
	require.NoError(t, err)
	require.NotEmpty(t, id)
	_, err = store.Get(id)
	require.Equal(t, sql.ErrNoRows, err)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	 // настройте подключение к БД
	 db, err := sql.Open("sqlite", "tracker.db")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer db.Close()
	store := NewParcelStore(db)

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	p := Parcel {
   		Client:    2,
    	Status:    ParcelStatusRegistered,
    	Address:   "testAdress",
    	Created_at: time.Now().UTC().Format(time.RFC3339),
	}
	var id int
	
	id, err = store.Add(p)
    require.NoError(t, err)
	require.NotEmpty(t, id)
	
	// set address
	// обновите адрес, убедитесь в отсутствии ошибки
	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
    require.NoError(t, err)

	// check
	// получите добавленную посылку и убедитесь, что адрес обновился
	var  returnParc Parcel
    returnParc, err = store.Get(id)
	require.NoError(t, err)
	require.NotEmpty(t, returnParc)
	assert.Equal(t, newAddress, returnParc.Address, "Не совпадает")
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	// настройте подключение к БД

	db, err := sql.Open("sqlite", "tracker.db")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer db.Close()
	store := NewParcelStore(db)

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
p := Parcel {
   		Client:    3,
    	Status:    ParcelStatusRegistered,
    	Address:   "testAdress",
    	Created_at: time.Now().UTC().Format(time.RFC3339),
	}
	id, err := store.Add(p)
    require.NoError(t, err)
	require.NotEmpty(t, id)

	// set status
	// обновите статус, убедитесь в отсутствии ошибки
	newStatus := ParcelStatusSent
	err = store.SetStatus(id, newStatus)
	require.NoError(t, err)

	// check
	// получите добавленную посылку и убедитесь, что статус обновился
	var  returnStatus Parcel
    returnStatus, err = store.Get(id)
	require.NoError(t, err)
	require.NotEmpty(t, returnStatus)
	assert.Equal(t, newStatus, returnStatus.Status, "Не совпадает")
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	// настройте подключение к БД
	db, err := sql.Open("sqlite", "tracker.db")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer db.Close() 
	store := NewParcelStore(db)
	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	// задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	var id int
	for i := 0; i < len(parcels); i++ {
		// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
		id, err = store.Add(parcels[i])
    require.NoError(t, err)
	require.NotEmpty(t, id)
	require.NotZero(t, id)
		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = id

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id] = parcels[i]
	}

	// get by client
	// получите список посылок по идентификатору клиента, сохранённого в переменной client
	storedParcels, err := store.GetByClient(client)
	// убедитесь в отсутствии ошибки
	require.NoError(t, err)
	// убедитесь, что количество полученных посылок совпадает с количеством добавленных
	require.Equal(t, len(parcels), len(storedParcels))
	// check
	for _, parcel := range storedParcels {
		// в parcelMap лежат добавленные посылки, ключ - идентификатор посылки, значение - сама посылка
		// убедитесь, что все посылки из storedParcels есть в parcelMap
		expectedParcel, found := parcelMap[parcel.Number]
		require.True(t, found) 
       
		// убедитесь, что значения полей полученных посылок заполнены верно
		assert.Equal(t, expectedParcel.Client, parcel.Client)
        assert.Equal(t, expectedParcel.Status, parcel.Status)
        assert.Equal(t, expectedParcel.Address, parcel.Address)
		pTime, _ := time.Parse(time.RFC3339, expectedParcel.Created_at)
		rTime, _ := time.Parse(time.RFC3339, parcel.Created_at)
		assert.WithinDuration(t, pTime, rTime, time.Second)
    }
	}

