<?php
// checks that the patterns work in PHP's preg_match

$data = json_decode(file_get_contents('crawler-user-agents.json'), true);

$patterns = array();
foreach($data as $entry) {
  $patterns[] = $entry['pattern'];
  if (isset($entry['instances'])) {
    foreach($entry['instances'] as $ua_example) {
      if (!preg_match('/'.$entry['pattern'].'/', $ua_example)) {
        throw new Exception('pb with '.$entry['pattern']);
      }
    }
  }
}

// testing with a giant regexp
$regexp = implode('|', $patterns);
foreach($data as $entry) {
  if (isset($entry['instances'])) {
    foreach($entry['instances'] as $ua_example) {
      if (!preg_match('/'.$regexp.'/', $ua_example)) {
        throw new Exception('pb with '.$entry['pattern']);
      }
    }
  }
}

?>
