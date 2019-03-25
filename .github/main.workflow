workflow "CI" {
  on = "push"
  resolves = ["Test"]
}

action "Test" {
  uses = "cedrickring/golang-action@707c4349bae930df82b058f89592a591e55b3dfa"
}
