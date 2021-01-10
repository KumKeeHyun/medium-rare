package dao_test

import (
	"testing"

	"github.com/KumKeeHyun/medium-rare/user-service/dao"
	"github.com/KumKeeHyun/medium-rare/user-service/dao/memory"
	"github.com/KumKeeHyun/medium-rare/user-service/domain"
)

func TestMethod(t *testing.T) {
	ur := memory.NewMemoryUserRepository()

	// testFindByID(t, ur)
	testFindByEmail(t, ur)
}

func testFindByID(t *testing.T, ur dao.UserRepository) {
	u1, _ := ur.Save(domain.User{
		Email: "test1@test.com", Password: "testpw", Name: "testName1",
	})

	findUser, _ := ur.FindByID(u1.ID)

	if findUser != u1 {
		t.Fail()
	}
}

func testFindByEmail(t *testing.T, ur dao.UserRepository) {
	u1, _ := ur.Save(domain.User{
		Email: "test1@test.com", Password: "testpw", Name: "testName1",
	})

	findUser, _ := ur.FindByEmail("test1@test.com")

	if findUser != u1 {
		t.Fail()
	}
}
