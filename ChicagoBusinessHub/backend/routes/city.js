var express = require('express');
var router = express.Router();

var data = []
var pg = require('pg');
var conString = "pg://postgres:HaveFun@127.0.0.1:5432/chicago_bi";
var pgClient = new pg.Client(conString);
pgClient.connect();

async function getCityData(query) {
  const resp = await pgClient.query(query);
  lData = []
  for (i = 0; i < resp.rows.length; i++) {
    var cd = {
      "fips": i+1,
      "zipcode": resp.rows[i].zipcode,
      "Poverty rate": resp.rows[i].below_poverty_percent,
      "CCVI": 0,
      "count": 0
    };
    lData.push(cd);
  }
  data.push(lData);
  return data;
}

async function getCityData_2(query, x) {
  const resp = await pgClient.query(query);
  //lData = []
  
  for (let i = 0; i < resp.rows.length; i++) {
    for (let j = 0; j < x[0].length; j++) {
      if (resp.rows[i].zipcode === x[0][j]["zipcode"]) {
        x[0][j]["CCVI"] = resp.rows[i].ccviscore
      }
    }
  }
  //data.push(lData);
  return x;
}

async function getCityData_3(query, x) {
  const resp = await pgClient.query(query);
  //lData = []
  
  for (let i = 0; i < resp.rows.length; i++) {
    for (let j = 0; j < x[0].length; j++) {
      if (resp.rows[i].zipcode === x[0][j]["zipcode"]) {
        x[0][j]["count"] = resp.rows[i].count
      }
    }
  }
  // console.log(x[0]);
  //data.push(lData);
  return x;
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
    text: "select * from poverty",
    //values: [ccviscore]
  }

  const query_2 = {
    // give the query a unique name
    //name: 'fetch-count-'+day,
    text: "select distinct(zipcode),ccviscore from ccvi_final order by zipcode",
    //values: [ccviscore]
  }

  const query_3 = {
    // give the query a unique name
    //name: 'fetch-count-'+day,
    text: "select zipcode, COUNT(*) from building_permits GROUP BY zipcode",
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
  
  getCityData(query).then(x => { 
    //res.json(x);
    //console.log("results from getdata ");
    getCityData_2(query_2, x).then(x => { 
      //console.log("results from getdata ");
      getCityData_3(query_3, x).then(x => {
        res.json(x);
      });
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
