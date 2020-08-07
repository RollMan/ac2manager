
(function(){
  function fetch_upcoming(){
    const upcoming_race_div = document.querySelector('div#upcoming_race')

    fetch("/api/upcoming_race", {
      headers:{
        method: 'GET'
      }
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
          throw new Error("Login Failed.")
        }
        return response.json()
      })
    .then(response => {
      if(response.length == 0){
        upcoming_race_div.innerHTML = "<p>No races scheduled.</p>";
      }else{
        const startdate = parseISO(response.startdate)
        const time_zone = 'Asia/Tokyo'
        const time_jst = utcToZonedTime(startdate, time_zone)
        const time_jst_string = formatRFC3339(time_jst)
        response.startdate = time_jst_string
        let body = JSON.stringify(response);
        upcoming_race_div.innerHTML = "<p>" + body + "</p>"
      }
    });
  }
  fetch_upcoming();
})();
