var createError = require('http-errors');
var express = require('express');
var path = require('path');
var cookieParser = require('cookie-parser');
var logger = require('morgan');

var heatmapRouter = require('./routes/heatmap');
var alertsRouter = require('./routes/alerts');
var chartsRouter = require('./routes/charts');
var dailydataRouter = require('./routes/dailydata');
var weekChartsRouter = require('./routes/weeklycharts');
var neighRouter = require('./routes/neighbourhood');
var ohareRouter = require('./routes/oharetrips');
var alertdatarouter = require('./routes/alertdata');
var midwayRouter = require('./routes/midwaytrips');
var highCCVIRouter = require('./routes/taxitripshighccvi');
var businessRouter = require('./routes/businessdata');
var checkEligibilityRouter = require('./routes/checkEligibility');
var monthlyRouter = require('./routes/monthly');
var weeklyRouter = require('./routes/weekly');
var dailyRouter = require('./routes/daily');
var cityRouter = require('./routes/city');
var cors = require('cors');

var app = express();

// view engine setup
app.set('views', path.join(__dirname, 'views'));
app.set('view engine', 'jade');

app.use(logger('dev'));
app.use(express.json());
app.use(express.urlencoded({ extended: false }));
app.use(cookieParser());
app.use(express.static(path.join(__dirname, 'public')));

app.use(cors());
app.use('/heatmap', heatmapRouter);
app.use('/alerts', alertsRouter);
app.use('/charts', chartsRouter);
app.use('/coviddailydata',dailydataRouter);
app.use('/covidweeklychart', weekChartsRouter);
app.use('/neighbourhooddata', neighRouter);
app.use('/oharetrips', ohareRouter);
app.use('/alertdata', alertdatarouter);
app.use('/midwaytrips', midwayRouter);
app.use('/businessdata', businessRouter);
app.use('/taxiTripsHighCCVI', highCCVIRouter);
app.use('/checkEligibility', checkEligibilityRouter);
app.use('/monthly', monthlyRouter);
app.use('/weekly', weeklyRouter);
app.use('/daily', dailyRouter);
app.use('/city', cityRouter);

// catch 404 and forward to error handler
app.use(function(req, res, next) {
  next(createError(404));
});

// error handler
app.use(function(err, req, res, next) {
  // set locals, only providing error in development
  res.locals.message = err.message;
  res.locals.error = req.app.get('env') === 'development' ? err : {};

  // render the error page
  res.status(err.status || 500);
  res.render('error');
});

module.exports = app;
