(function(){

const submit_button = document.querySelector('button#login')
const userid_input = document.querySelector('input[name="userid"]')
const pw_input = document.querySelector('input[name="pw"]')
const loginresult_div = document.querySelector('div#login_result')

submit_button.addEventListener("click", function(){
  let login_request_header = new Headers();
  login_request_header.append('Content-Type', 'application/json')


  let userid = userid_input.value
  let pw = pw_input.value
  let login_request_body = JSON.stringify({'userid': userid, 'pwhash': pw})

  fetch("/api/login", {
    method: 'POST',
    headers: login_request_header,
    body: login_request_body
  })
    .then(response => {
      if (!response.ok) {
        let body = "";
        if (400 <= response.status && response.status < 500){
          body = "<p>Authentication Failed.</p>"
        }else if(500 <= response.status && response.status < 600){
          body = "<p>Internal Server Error.</p>"
        }else{
          body = "<p>Unknown Error: " + response.status + ".</p>"
        }
        loginresult_div.innerHTML = body;
        throw new Error("Login Failed.")
      }
      return response.json()
    })
    .then(response => {
      let jwt = response.jwt
      document.cookie = "jwt=" + jwt
      loginresult_div.innerHTML = "ok";
    });
});

})();
