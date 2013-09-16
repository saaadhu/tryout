var app = angular.module('tryout', ['ui.bootstrap']);

function TryoutCtrl ($http, $scope) {
    $scope.options = '';

    $scope.compile = function() {
        $http.post ("/compile",
                { "code" : $scope.code, "options" : $scope.options }).success (function (response) {
                  $scope.output = response.BuildOutput;
                  $scope.lss = response.Listing;
                });
    };
}
