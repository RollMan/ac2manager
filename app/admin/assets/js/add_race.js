(function(){
  const event_table_for_edit_div = document.querySelector('div#event_table_for_edit')
  const result_div = document.querySelector('div#event_add_result')
  const submit_button = document.querySelector('button#submit_race')

  let schema;

  fetch('/api/schema', {
    method: 'GET',
    headers: {'Content-Type': 'application/json'}
  })
    .then(response => {
      if (!response.ok) {
        let body = "";
        if (400 <= response.status && response.status < 500){
          body = "<p>Client Error.</p>"
        }else if(500 <= response.status && response.status < 600){
          body = "<p>Internal Server Error.</p>"
        }else{
          body = "<p>Unknown Error: " + response.status + ".</p>"
        }
        upcoming_race_div.innerHTML = body;
        throw new Error("Server Unavailable.")
      }
      return response.json()
    })
    .then(response => {
      schema = response;
    });
})();
