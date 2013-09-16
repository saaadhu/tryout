var app = angular.module('tryout', ['ui.bootstrap']);

function TryoutCtrl ($http, $scope) {
    $scope.options = '';
    $scope.status_color = '#FFFFFF';

    $scope.compile = function() {
        $scope.output = '';
        $scope.lss = '';
        $scope.status_color = '#FFFFFF';
        $http.post ("/compile",
                { "code" : $scope.code, "options" : $scope.options }).success (function (response) {
                  $scope.output = response.BuildOutput;
                  $scope.lss = response.Listing;
                  $scope.status_color = response.Success ? "#FFFFFF" : '#FFCCCC';
                });
    };
}
