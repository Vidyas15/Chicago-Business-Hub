var express = require('express');
var router = express.Router();

var data = []
var pg = require('pg');
var conString = "pg://postgres:HaveFun@127.0.0.1:5432/chicago_bi";
var pgClient = new pg.Client(conString);
pgClient.connect();

async function getDailyTrafficData(query) {
  const resp = await pgClient.query(query);
  lData1 = []
  lData2 = []
  vData = []
  //console.log('entering for loop');
  //console.log(resp.rows.length);
  for (i = 0; i < resp.rows.length; i++) {
    //console.log(resp.rows[i]);
    // var cd = {
    //   "neighborhood": resp.rows[i].neighborhood,
    //   "ccviscore": resp.rows[i].ccviscore,
    // };
    //console.log(typeof String(resp.rows[i].Date))
    //console.log(typeof resp.rows[i].Date);
    //if(resp.rows[i].Date<"2021-11-02")
      lData1.push(resp.rows[i].date1);
    //else
      //lData2.push(resp.rows[i].Date);
    vData.push(resp.rows[i].traffic);
    //lData.push(cd);
  }
  console.log(lData1);
  console.log(lData2);
  data.push(lData1);
  data.push(lData2);
  data.push(vData);
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
    text: "select * from mo_oh",
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
  
  getDailyTrafficData(query).then(x => { 
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