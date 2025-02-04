package service

import (
	"hoge-api/auth"
	"hoge-api/model"
	"hoge-api/security"

	"github.com/google/uuid"
)

type UserService struct{} // コントローラ側からサービスを実体として使うため。この構造体にレシーバ機能でメソッドを紐づける。

func (s *UserService) RegisterUser(bUser model.User) (string, error) {
	// user_uuidを生成
	userId, err := uuid.NewRandom() // 新しいuuidの生成
	if err != nil {
		return "", err
	}
	bUser.UserUuid = userId.String() // 設定

	// パスワードをハッシュ化
	hashed, err := security.HashingByEncrypt(bUser.Password)
	if err != nil {
		return "", NewErr(ErrTypeHashingPassFailed, err.Error()) // controller側でエラーを判別するために、ENUMで管理されたエラーを返す
	}
	// ハッシュ(:バイト配列)化されたパスワードを文字列にsh知恵構造体に戻す
	bUser.Password = string(hashed)

	// 構造体をレコード登録処理に投げる
	_, err = model.CreateUser(bUser) // 第一返り血は登録成功したレコード数
	if err != nil {
		return "", err
	}

	// 登録が成功したらトークンを生成する
	token, err := auth.GenerateToken(bUser.UserUuid) // トークンを取得
	if err != nil {
		return "", NewErr(ErrTypeGenTokenFailed, err.Error())
	}
	return token, nil
}
