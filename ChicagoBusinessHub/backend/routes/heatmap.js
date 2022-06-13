var express = require('express');
var router = express.Router();

var data = []
var pg = require('pg');
var conString = "pg://postgres:HaveFun@127.0.0.1:5432/chicago_bi";
var pgClient = new pg.Client(conString);
pgClient.connect();

async function getCovidHeatMapData(query) {
  const resp = await pgClient.query(query);
  lData = []
  //console.log('entering for loop');
  for (i = 0; i < resp.rows.length; i++) {  
    //console.log(resp.rows[i]);
    var cd = {
      "ccvi": resp.rows[i].ccviscore,
      "latitude": resp.rows[i].latitude,
      "longitude": resp.rows[i].longitude
    };
    lData.push(cd);
  }
  //console.log(lData);
  data.push(lData);
  return data;
}

router.get('/', function(req, res, next) {
  const query = {
    text: 'select * from ccvi_final'
  }

  getCovidHeatMapData(query).then(x => { 
    //console.log("results from getdata ");
    res.json(x);
  });

  //res.render('index', { title: 'Express' });
});

module.exports = router;
