INSERT INTO words (word, part_of_speech, definition, example, difficulty) VALUES 
('abandon', 'verb', 'To leave completely and finally', 'They had to abandon the sinking ship', 'easy'),
('benefit', 'noun', 'An advantage or profit gained from something', 'One benefit of exercise is better sleep', 'easy'),
('capable', 'adjective', 'Having the ability or quality needed to do something', 'She is capable of solving complex problems', 'easy'),
('despite', 'preposition', 'Without being affected by', 'He went running despite the heavy rain', 'easy'),
('eager', 'adjective', 'Strongly wanting to do or have something', 'The students were eager to learn new skills', 'easy')
ON CONFLICT (word) DO NOTHING;