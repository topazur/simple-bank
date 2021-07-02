package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "github.com/topaz-h/go-simple-bank/db/mock"
	db "github.com/topaz-h/go-simple-bank/db/sqlc"
	"github.com/topaz-h/go-simple-bank/util"
)

/// GetAccount api测试用例 `单个情况`
func TestGetAccountAPI(t *testing.T) {
	// 表关联即可
	user, _ := randomUser(t)
	account := randomAccount(user.Username)

	// 传入 *testing.T 生成 *gomock.Controller
	ctrl := gomock.NewController(t)

	// defer执行：将检查是否所有预期被调用的方法都被调用
	// defer ctrl.Finish() // 新版本不再需要 (当你将 *testing.T 传递给 NewController 函数)

	store := mockdb.NewMockStore(ctrl)

	/// build stubs: 我希望调用Store的 GetAccount 函数,具有任何上下文和此特定帐户 ID 参数。
	// 第一个上下文参数 `ctx context.Context` 可以是任何值，所以我们使用gomock.Any() 匹配器
	// Times(1) 意味着我们希望这个函数被精确调用 1 次
	// 使用 Return 函数告诉 gomock 返回一些特定的值；与被调用函数的返回值匹配
	store.EXPECT().
		GetAccount(gomock.Any(), gomock.Eq(account.ID)).
		Times(1).
		Return(account, nil)

	// 启动测试 HTTP 服务器并发送请求
	server, _ := NewServer(store)

	// 使用 httptest 包的 Recorder 功能记录 API 请求的响应
	recorder := httptest.NewRecorder()

	// http.NewRequest
	url := fmt.Sprintf("/account/%d", account.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	// 将通过服务器路由器发送我们的 API 请求,并将其响应记录在记录器(recorder)中。
	server.router.ServeHTTP(recorder, request)

	// check response (Code or Body)
	require.Equal(t, http.StatusOK, recorder.Code)
	requireBodyMatchAccount(t, recorder.Body, account)
}

/// GetAccount api测试用例列表 `所有情况`
func TestGetAccountAPIAll(t *testing.T) {
	user, _ := randomUser(t)
	account := randomAccount(user.Username)

	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name:      "1-OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		// TODO: add more cases...
		{
			name:      "2-NotFound",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "3-InternalError",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "4-InvalidID",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			// defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store) // Add

			server, _ := NewServer(store)
			recorder := httptest.NewRecorder()

			// NOTICE: tc.accountID 当前测试用例使用的参数
			url := fmt.Sprintf("/account/%d", tc.accountID) // Modify
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder) // Add
		})
	}
}

/// 辅助函数 - 生成一个随机帐户。
func randomAccount(owner string) db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    owner,
		Balance:  util.RandomAmountMoney(),
		Currency: util.RandomCurrency(),
	}
}

/// 辅助函数 - 验证响应主体是否与account匹配。
func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	// 首先我们调用 ioutil.ReadAll() 从响应体中读取所有数据,保存为json格式
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)

	require.Equal(t, account, gotAccount)
}
