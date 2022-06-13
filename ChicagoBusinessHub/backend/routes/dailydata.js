var express = require('express');
var router = express.Router();

var data = []
var pg = require('pg');
var conString = "pg://postgres:HaveFun@127.0.0.1:5432/chicago_bi";
var pgClient = new pg.Client(conString);
pgClient.connect();

async function getWeeklyData(query) {
  const resp = await pgClient.query(query);
  lData = []
  for (i = 0; i < resp.rows.length; i++) {
    //console.log(resp.rows[i]);
    var cd = {
      "index": i+1,
      "week": resp.rows[i].week,
      "community_area": resp.rows[i].communityareas,
      "cases_weekly": resp.rows[i].cases_weekly,
      "cases_cumulative": resp.rows[i].cases_cumulative,
      "deaths_weekly": resp.rows[i].deaths_weekly,
      "deaths_cumulative": resp.rows[i].deaths_cumulative,
      "population": resp.rows[i].population,
      "tests_weekly": resp.rows[i].tests_weekly,
    };
    lData.push(cd);
  }
  //console.log(lData);
  data.push(lData);
  return data;
}

/* GET users listing. */
router.get('/', function(req, res, next) {
  //console.log('function entry');
  /*let jsonResponse = {
    "handsets":[
      {title:'Raw Data', cols:2, rows:1},
      {title:'Charts', cols:2, rows:1},
      {title:'Heat Map', cols:2, rows:1},
      {title:'Alert system', cols:2, rows:1}
      ],
    "web":[
      {title:'Raw Data', cols:2, rows:1},
          {title:'Charts', cols:2, rows:1},
          {title:'Heat Map', cols:2, rows:1},
          {title:'Alert system', cols:2, rows:1}
      ]
  };*/
  


  //var qry = `SELECT * from covid_ccvi`;
  const query = {
    // give the query a unique name
    //name: 'fetch-count-'+day,
    text: "select * from covid_zip where week='11/20/2021' order by cases_weekly desc",
    //values: [ccviscore]
  }

  // (async () => {
  //   console.log(await getData(query))
  // })();
  // res.send('OK');
  // let ans = function() {
  //   getData(query).then(function(value){
  //     console.log(value);
  //   });
  //   console.log(data);
  //   res.send('OK')
  // }

  //resp = getData(query);
  
  getWeeklyData(query).then(x => { 
    //console.log("results from getdata ");
    res.json(x);
  });
  // setTimeout(function(){ res.json(resp); }, 5000);
  //res.json(resp);
  //res.send('OK');
  // getData(query).then(function (response) {
  //   var hits = response;
  //   res.json({'data found': 'Successfully Retrieved'});
  // });

  //console.log(resp);
  // for (i = 1; i < 10; i++) {
  //   console.log(response);
  // }

  //console.log(data);
  //res.json(data);
});

module.exports = router;
