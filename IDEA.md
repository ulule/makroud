# Query Builder

$flights = App\Flight::where('active', 1)
               ->orderBy('name', 'desc')
               ->take(10)
               ->get();

# Chuck

Flight::chunk(200, function ($flights) {
    foreach ($flights as $flight) {
        //
    }
});

# Cursor

foreach (Flight::where('foo', 'bar')->cursor() as $flight) {
    //
}

# Find

App\Flight::find([1, 2, 3]);

# Updates

App\Flight::where('active', 1)
          ->where('destination', 'San Diego')
          ->update(['delayed' => 1]);
