/**
 * @Author fzh
 * @Date 2021/2/21
 */
package userrpc

import "testing"

var client = NewClient()

func TestClient_GetUser(t *testing.T) {
	user, err := client.GetUser(1)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%#v\n", user)
}
