-- Remove seed test data (cascades to workouts, exercises, templates, schedule via FK)
DELETE FROM users WHERE login IN ('alexei_ivanov','maria_petrova','dmitry_smirnov','elena_kozlova');
