package model

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"time"
	"wangqingang/cunxun/db"
	"wangqingang/cunxun/test"
)

func TestCreateCarpooling(t *testing.T) {
	test.InitTestCaseEnv(t)
	assert := assert.New(t)

	u := &User{
		Phone:          test.GenFakePhone(),
		NickName:       test.GenRandString(),
		HashedPassword: test.GenRandString(),
		PasswordLevel:  test.GenRandInt(5),
		RegisterSource: test.TestWebSource,
		Avatar:         test.GenRandString(),
	}

	u, err := CreateUser(db.Mysql, u)
	assert.Nil(err)
	assert.NotNil(u)

	c := &Carpooling{
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		CreaterUid:  u.ID,
		UpdaterUid:  u.ID,
		FromCity:    test.GenRandString(),
		ToCity:      test.GenRandString(),
		DepartTIme:  time.Now(),
		PeopleCount: test.GenRandInt(5),
		Contact:     test.GenRandString(),
		Status:      test.GenRandInt(1),
		Remark:      test.GenRandString() + test.GenRandString(),
	}

	c, err = CreateCarpooling(db.Mysql, c)
	assert.Nil(err)
	assert.NotNil(c)

	cd, err := GetCarpoolingByID(db.Mysql, c.ID)
	assert.Nil(err)
	assert.NotNil(cd)
}

func TestGetCarpoolingList(t *testing.T) {
	test.InitTestCaseEnv(t)
	assert := assert.New(t)

	u := &User{
		Phone:          test.GenFakePhone(),
		NickName:       test.GenRandString(),
		HashedPassword: test.GenRandString(),
		PasswordLevel:  test.GenRandInt(5),
		RegisterSource: test.TestWebSource,
		Avatar:         test.GenRandString(),
	}

	u, err := CreateUser(db.Mysql, u)
	assert.Nil(err)
	assert.NotNil(u)

	var cs []*Carpooling
	for i := 0; i < 10; i++ {
		c := &Carpooling{
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			CreaterUid:  u.ID,
			UpdaterUid:  u.ID,
			FromCity:    test.GenRandString(),
			ToCity:      test.GenRandString(),
			DepartTIme:  time.Now(),
			PeopleCount: test.GenRandInt(5),
			Contact:     test.GenRandString(),
			Status:      test.GenRandInt(1),
			Remark:      test.GenRandString() + test.GenRandString(),
		}
		_, err := CreateCarpooling(db.Mysql, c)
		assert.Nil(err)
		cs = append(cs, c)
	}

	items, isOver, err := GetCarpoolingList(db.Mysql, map[string]interface{}{}, OrderByIgnore, 20, 1)
	assert.Nil(err)
	assert.NotNil(items)
	assert.True(isOver)
	assert.Equal(10, len(items))
}

func TestUpdateCarpoolingById(t *testing.T) {
	test.InitTestCaseEnv(t)
	assert := assert.New(t)

	u := &User{
		Phone:          test.GenFakePhone(),
		NickName:       test.GenRandString(),
		HashedPassword: test.GenRandString(),
		PasswordLevel:  test.GenRandInt(5),
		RegisterSource: test.TestWebSource,
		Avatar:         test.GenRandString(),
	}

	u, err := CreateUser(db.Mysql, u)
	assert.Nil(err)
	assert.NotNil(u)

	c := &Carpooling{
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		CreaterUid:  u.ID,
		UpdaterUid:  u.ID,
		FromCity:    test.GenRandString(),
		ToCity:      test.GenRandString(),
		DepartTIme:  time.Now(),
		PeopleCount: test.GenRandInt(5),
		Contact:     test.GenRandString(),
		Status:      test.GenRandInt(1),
		Remark:      test.GenRandString() + test.GenRandString(),
	}

	c, err = CreateCarpooling(db.Mysql, c)
	assert.Nil(err)
	assert.NotNil(c)

	cd, err := GetCarpoolingByID(db.Mysql, c.ID)
	assert.Nil(err)
	assert.NotNil(cd)

	newFromCity := test.GenRandString()
	count, err := UpdateCarpoolingById(db.Mysql, c.ID, &Carpooling{FromCity: newFromCity})
	assert.Nil(err)
	assert.NotZero(count)

	cd, err = GetCarpoolingByID(db.Mysql, c.ID)
	assert.Nil(err)
	assert.NotNil(cd)
	assert.Equal(newFromCity, cd.FromCity)
}

func TestDeleteCarpoolingById(t *testing.T) {
	test.InitTestCaseEnv(t)
	assert := assert.New(t)

	u := &User{
		Phone:          test.GenFakePhone(),
		NickName:       test.GenRandString(),
		HashedPassword: test.GenRandString(),
		PasswordLevel:  test.GenRandInt(5),
		RegisterSource: test.TestWebSource,
		Avatar:         test.GenRandString(),
	}

	u, err := CreateUser(db.Mysql, u)
	assert.Nil(err)
	assert.NotNil(u)

	c := &Carpooling{
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		CreaterUid:  u.ID,
		UpdaterUid:  u.ID,
		FromCity:    test.GenRandString(),
		ToCity:      test.GenRandString(),
		DepartTIme:  time.Now(),
		PeopleCount: test.GenRandInt(5),
		Contact:     test.GenRandString(),
		Status:      test.GenRandInt(1),
		Remark:      test.GenRandString() + test.GenRandString(),
	}

	c, err = CreateCarpooling(db.Mysql, c)
	assert.Nil(err)
	assert.NotNil(c)

	cd, err := GetCarpoolingByID(db.Mysql, c.ID)
	assert.Nil(err)
	assert.NotNil(cd)

	count, err := DeleteCarpoolingById(db.Mysql, c.ID)
	assert.Nil(err)
	assert.NotZero(count)

	cd, err = GetCarpoolingByID(db.Mysql, c.ID)
	assert.Nil(err)
	assert.Nil(cd)
}
