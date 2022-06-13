var express = require('express');
var router = express.Router();

var data = []
var pg = require('pg');
var conString = "pg://postgres:HaveFun@127.0.0.1:5432/chicago_bi";
var pgClient = new pg.Client(conString);
pgClient.connect();

async function getBusinessData(query, key) {
  const resp = await pgClient.query(query);
  lData = []
  var val;
  for (i = 0; i < resp.rows.length; i++) {
    console.log(resp.rows[i]);
    if (key == resp.rows[i].zipcode) {
      val = true;
    } else {
      val = false;
    }
    var cd = {
      "index": i+1,
      "Permit No": resp.rows[i].permit_no,
      "Permit Type": resp.rows[i].permit_type,
      "Application Date": resp.rows[i].application_date,
      "Issue Date": resp.rows[i].issue_date,
      "Zipcode": resp.rows[i].zipcode,
      "Community Area": resp.rows[i].community_area,
      "No_of_Applications": resp.rows[i].no_of_applications
    };
    lData.push(cd);
  }
  console.log(val);
  data.push(lData);
  return data;
}

/* GET users listing. */
router.get('/', function(req, res, next) {
  //console.log('function entry');
  console.log(req.query.searchKey);
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
    text: "select distinct(b.permit_no),b.permit_type,b.application_date,b.issue_date,b.zipcode,b.community_area,(count(b.zipcode) OVER ( partition by b.zipcode))/2 as no_of_applications from building_permits b inner join poverty p on b.zipcode=p.zipcode and b.zipcode='60620' and p.per_capita_income < 30000 order by no_of_applications",
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
  
  getBusinessData(query, req.query.searchKey).then(x => { 
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
