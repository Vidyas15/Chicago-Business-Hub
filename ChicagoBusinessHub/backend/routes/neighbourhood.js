var express = require('express');
var router = express.Router();

var data = []
var pg = require('pg');
var conString = "pg://postgres:HaveFun@127.0.0.1:5432/chicago_bi";
var pgClient = new pg.Client(conString);
pgClient.connect();

async function getNeighData(query) {
  const resp = await pgClient.query(query);
  lData = []
  for (i = 0; i < 5; i++) {
    //console.log(resp.rows[i]);
    var cd = {
      "index": i+1,
      "Community area": resp.rows[i].community_area,
      "Poverty rate": resp.rows[i].below_poverty_percent,
    };
    lData.push(cd);
  }
  //console.log(lData);
  data.push(lData);
  return data;
}

async function getNeighData_2(query) {
  const resp = await pgClient.query(query);
  lData = []
  for (i = 0; i < 5; i++) {
    //console.log(resp.rows[i]);
    var cd = {
      "index": i+1,
      "Community Area": resp.rows[i].community_area,
      "Unemployment rate": resp.rows[i].percent_unemployed,
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
    text: "select * from poverty order by below_poverty_percent desc",
    //values: [ccviscore]
  }

  const query_2 = {
    // give the query a unique name
    //name: 'fetch-count-'+day,
    text: "select * from poverty order by percent_unemployed desc",
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
  
  getNeighData(query).then(x => { 
    //console.log("results from getdata ");
    getNeighData_2(query_2).then(x => { 
      //console.log("results from getdata ");
      res.json(x);
    });
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
