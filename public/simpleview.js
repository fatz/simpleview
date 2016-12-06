var app = angular.module("SimpleView", []);

app.controller("OverviewCtrl", function($scope, $http, $timeout, $rootScope) {
  $scope.loadData = function (){
    $http.get('v1/overview').
      success(function(data, status, headers, config) {
        $rootScope.bodyClass = 'bodynormal';
        $scope.error = [];
        $scope.error['data'] = []
        $scope.downs = data['hosts'];
        criticals = [];
        warns = [];

        if (data['services']) {
          data['services'].forEach(function(item, index){
            if (item['attrs']['last_hard_state'] == 2) {
              criticals.push(item);
            } else if (item['attrs']['last_hard_state'] == 1) {
              warns.push(item);
            }
          })
        }

        $scope.criticals = criticals;
        $scope.warns = warns;
      }).
      error(function(data, status, headers, config) {
          $rootScope.bodyClass = 'bodyerror';
          $scope.error = [data];
      });
    }

    $scope.intervalFunction = function(){
   $timeout(function() {
     $scope.loadData();
     $scope.intervalFunction();
   }, 10000)
 };
    $scope.loadData();
    $scope.intervalFunction();

});
