-- Seed: 4 test clients, workouts, schedule entries, templates.
-- Password for all test users: 1234  (bcrypt cost=10)
-- Trainer: macaretchurilov (id=1, role=admin)

BEGIN;

-- ── Users ─────────────────────────────────────────────────────────────────────
INSERT INTO users (login, email, password_hash, role) VALUES
  ('alexei_ivanov',  'alexei@gym.local',  '$2a$10$xlm2Y4YiWsEnf4tuEclIyuZyhkSpFvXwpRWD/IFYVzHpZ0orQEElm', 'user'),
  ('maria_petrova',  'maria@gym.local',   '$2a$10$xlm2Y4YiWsEnf4tuEclIyuZyhkSpFvXwpRWD/IFYVzHpZ0orQEElm', 'user'),
  ('dmitry_smirnov', 'dmitry@gym.local',  '$2a$10$xlm2Y4YiWsEnf4tuEclIyuZyhkSpFvXwpRWD/IFYVzHpZ0orQEElm', 'user'),
  ('elena_kozlova',  'elena@gym.local',   '$2a$10$xlm2Y4YiWsEnf4tuEclIyuZyhkSpFvXwpRWD/IFYVzHpZ0orQEElm', 'user')
ON CONFLICT (login) DO NOTHING;

-- ── Workouts — alexei_ivanov (powerlifting focus) ─────────────────────────────
-- 1. Грудь + Трицепс — 07.04
WITH w AS (
  INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Грудь + Трицепс', '2026-04-07', 75, 'Хороший день, веса шли легко'
  FROM users WHERE login = 'alexei_ivanov' RETURNING id
)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.exercise_id, e.sets, e.reps, e.weight_kg
FROM w, (VALUES (1,4,8,100.0),(2,4,8,32.0),(3,3,10,80.0),(26,4,8,70.0),(28,3,12,27.5)) AS e(exercise_id,sets,reps,weight_kg);

-- 2. Спина + Бицепс — 14.04
WITH w AS (
  INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Спина + Бицепс', '2026-04-14', 80, 'Добавил вес на тягу штанги'
  FROM users WHERE login = 'alexei_ivanov' RETURNING id
)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.exercise_id, e.sets, e.reps, e.weight_kg
FROM w, (VALUES (9,4,8,90.0),(10,3,10,28.0),(11,3,12,55.0),(21,4,10,45.0),(23,3,12,16.0)) AS e(exercise_id,sets,reps,weight_kg);

-- 3. Ноги — 21.04
WITH w AS (
  INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'День ног', '2026-04-21', 90, 'Приседания 110 кг — личник!'
  FROM users WHERE login = 'alexei_ivanov' RETURNING id
)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.exercise_id, e.sets, e.reps, e.weight_kg
FROM w, (VALUES (31,5,5,110.0),(32,4,12,150.0),(38,3,12,60.0),(34,3,15,30.0),(37,4,20,0.0)) AS e(exercise_id,sets,reps,weight_kg);

-- 4. Плечи — 28.04
WITH w AS (
  INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Плечи', '2026-04-28', 65, ''
  FROM users WHERE login = 'alexei_ivanov' RETURNING id
)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.exercise_id, e.sets, e.reps, e.weight_kg
FROM w, (VALUES (15,4,8,60.0),(16,3,10,22.0),(17,4,15,12.0),(18,3,12,40.0),(20,3,10,20.0)) AS e(exercise_id,sets,reps,weight_kg);

-- 5. Грудь + Трицепс — 05.05
WITH w AS (
  INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Грудь + Трицепс', '2026-05-05', 80, 'Жим 102.5 кг — прогресс'
  FROM users WHERE login = 'alexei_ivanov' RETURNING id
)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.exercise_id, e.sets, e.reps, e.weight_kg
FROM w, (VALUES (1,4,8,102.5),(2,4,8,34.0),(3,3,10,82.5),(27,4,10,40.0),(28,3,12,30.0)) AS e(exercise_id,sets,reps,weight_kg);

-- 6. Спина + Бицепс — 14.05
WITH w AS (
  INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Спина + Бицепс', '2026-05-14', 75, ''
  FROM users WHERE login = 'alexei_ivanov' RETURNING id
)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.exercise_id, e.sets, e.reps, e.weight_kg
FROM w, (VALUES (9,4,8,92.5),(11,3,12,57.5),(12,3,12,50.0),(21,4,10,47.5),(22,3,12,18.0)) AS e(exercise_id,sets,reps,weight_kg);

-- 7. Ноги — 26.05
WITH w AS (
  INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'День ног', '2026-05-26', 85, 'Устал, но все подходы выполнил'
  FROM users WHERE login = 'alexei_ivanov' RETURNING id
)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.exercise_id, e.sets, e.reps, e.weight_kg
FROM w, (VALUES (31,5,5,112.5),(32,4,12,155.0),(14,3,5,120.0),(34,3,15,32.5),(35,3,15,35.0)) AS e(exercise_id,sets,reps,weight_kg);

-- ── Workouts — maria_petrova (full-body / cardio) ─────────────────────────────
-- 1. Фулл-боди — 09.04
WITH w AS (
  INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Фулл-боди', '2026-04-09', 60, 'Первая тренировка после паузы'
  FROM users WHERE login = 'maria_petrova' RETURNING id
)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.exercise_id, e.sets, e.reps, e.weight_kg
FROM w, (VALUES (5,3,15,0.0),(33,3,12,8.0),(17,3,15,6.0),(39,3,20,0.0),(45,1,30,0.0)) AS e(exercise_id,sets,reps,weight_kg);

-- 2. Кардио + Пресс — 16.04
WITH w AS (
  INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Кардио + Пресс', '2026-04-16', 45, 'Кардио 20 мин + пресс'
  FROM users WHERE login = 'maria_petrova' RETURNING id
)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.exercise_id, e.sets, e.reps, e.weight_kg
FROM w, (VALUES (45,1,20,0.0),(46,1,15,0.0),(39,4,20,0.0),(40,3,60,0.0),(41,3,15,0.0)) AS e(exercise_id,sets,reps,weight_kg);

-- 3. Фулл-боди — 23.04
WITH w AS (
  INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Фулл-боди', '2026-04-23', 65, 'Добавила вес на выпады'
  FROM users WHERE login = 'maria_petrova' RETURNING id
)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.exercise_id, e.sets, e.reps, e.weight_kg
FROM w, (VALUES (5,3,15,0.0),(33,3,12,10.0),(2,3,12,14.0),(17,3,15,7.0),(39,3,20,0.0)) AS e(exercise_id,sets,reps,weight_kg);

-- 4. Верх тела — 30.04
WITH w AS (
  INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Верх тела', '2026-04-30', 70, 'Жим гантелей лёжа — прогресс!'
  FROM users WHERE login = 'maria_petrova' RETURNING id
)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.exercise_id, e.sets, e.reps, e.weight_kg
FROM w, (VALUES (2,3,12,16.0),(10,3,12,14.0),(17,4,15,8.0),(22,3,12,8.0),(40,4,45,0.0)) AS e(exercise_id,sets,reps,weight_kg);

-- 5. Кардио + Пресс — 07.05
WITH w AS (
  INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Кардио + Пресс', '2026-05-07', 50, ''
  FROM users WHERE login = 'maria_petrova' RETURNING id
)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.exercise_id, e.sets, e.reps, e.weight_kg
FROM w, (VALUES (47,1,10,0.0),(45,1,25,0.0),(39,4,20,0.0),(41,3,15,0.0),(44,3,20,0.0)) AS e(exercise_id,sets,reps,weight_kg);

-- 6. Фулл-боди — 19.05
WITH w AS (
  INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Фулл-боди', '2026-05-19', 60, 'Стабильно, без рекордов'
  FROM users WHERE login = 'maria_petrova' RETURNING id
)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.exercise_id, e.sets, e.reps, e.weight_kg
FROM w, (VALUES (5,3,15,0.0),(33,3,12,10.0),(2,3,12,16.0),(39,4,20,0.0),(40,3,60,0.0)) AS e(exercise_id,sets,reps,weight_kg);

-- ── Workouts — dmitry_smirnov (heavy compound lifts) ─────────────────────────
-- 1. Приседания + Ноги — 08.04
WITH w AS (
  INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Ноги — тяжёлый день', '2026-04-08', 95, 'Присед 130 кг × 5'
  FROM users WHERE login = 'dmitry_smirnov' RETURNING id
)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.exercise_id, e.sets, e.reps, e.weight_kg
FROM w, (VALUES (31,5,5,130.0),(14,3,5,140.0),(32,4,10,180.0),(34,3,12,45.0),(37,5,20,0.0)) AS e(exercise_id,sets,reps,weight_kg);

-- 2. Жим + Грудь — 17.04
WITH w AS (
  INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Грудь — тяжёлый день', '2026-04-17', 85, 'Жим 125 кг — новый максимум'
  FROM users WHERE login = 'dmitry_smirnov' RETURNING id
)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.exercise_id, e.sets, e.reps, e.weight_kg
FROM w, (VALUES (1,5,5,125.0),(3,3,8,95.0),(2,3,10,40.0),(29,3,10,0.0),(26,3,8,80.0)) AS e(exercise_id,sets,reps,weight_kg);

-- 3. Становая + Спина — 24.04
WITH w AS (
  INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Спина — тяжёлый день', '2026-04-24', 90, 'Становая 160 кг — чисто'
  FROM users WHERE login = 'dmitry_smirnov' RETURNING id
)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.exercise_id, e.sets, e.reps, e.weight_kg
FROM w, (VALUES (14,5,5,160.0),(9,4,8,110.0),(8,3,8,0.0),(10,3,10,40.0),(13,3,15,50.0)) AS e(exercise_id,sets,reps,weight_kg);

-- 4. Ноги — 06.05
WITH w AS (
  INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Ноги — тяжёлый день', '2026-05-06', 90, ''
  FROM users WHERE login = 'dmitry_smirnov' RETURNING id
)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.exercise_id, e.sets, e.reps, e.weight_kg
FROM w, (VALUES (31,5,5,132.5),(14,3,5,142.5),(32,4,10,185.0),(38,3,10,80.0),(37,5,20,0.0)) AS e(exercise_id,sets,reps,weight_kg);

-- 5. Грудь + Плечи — 15.05
WITH w AS (
  INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Грудь + Плечи', '2026-05-15', 80, 'Армейский жим идёт хорошо'
  FROM users WHERE login = 'dmitry_smirnov' RETURNING id
)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.exercise_id, e.sets, e.reps, e.weight_kg
FROM w, (VALUES (1,4,8,115.0),(15,4,5,80.0),(3,3,10,90.0),(17,4,12,20.0),(18,3,10,55.0)) AS e(exercise_id,sets,reps,weight_kg);

-- 6. Спина — 22.05
WITH w AS (
  INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Спина', '2026-05-22', 75, ''
  FROM users WHERE login = 'dmitry_smirnov' RETURNING id
)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.exercise_id, e.sets, e.reps, e.weight_kg
FROM w, (VALUES (14,3,5,145.0),(9,4,8,112.5),(11,3,12,62.5),(12,3,12,55.0),(13,3,15,55.0)) AS e(exercise_id,sets,reps,weight_kg);

-- 7. Ноги — 02.06
WITH w AS (
  INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Ноги — тяжёлый день', '2026-06-02', 95, 'Присед 135 кг — личник!'
  FROM users WHERE login = 'dmitry_smirnov' RETURNING id
)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.exercise_id, e.sets, e.reps, e.weight_kg
FROM w, (VALUES (31,5,5,135.0),(14,4,5,147.5),(32,4,10,190.0),(34,3,12,47.5),(37,5,20,0.0)) AS e(exercise_id,sets,reps,weight_kg);

-- ── Workouts — elena_kozlova (light / toning) ─────────────────────────────────
-- 1. Верх тела — 10.04
WITH w AS (
  INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Верх тела', '2026-04-10', 55, 'Начинаю заново, лёгкие веса'
  FROM users WHERE login = 'elena_kozlova' RETURNING id
)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.exercise_id, e.sets, e.reps, e.weight_kg
FROM w, (VALUES (5,3,12,0.0),(2,3,12,8.0),(17,3,15,5.0),(22,3,12,5.0),(40,3,30,0.0)) AS e(exercise_id,sets,reps,weight_kg);

-- 2. Пресс + Кардио — 18.04
WITH w AS (
  INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Пресс + Кардио', '2026-04-18', 40, ''
  FROM users WHERE login = 'elena_kozlova' RETURNING id
)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.exercise_id, e.sets, e.reps, e.weight_kg
FROM w, (VALUES (49,1,25,0.0),(39,3,20,0.0),(40,3,30,0.0),(41,3,12,0.0),(42,3,20,0.0)) AS e(exercise_id,sets,reps,weight_kg);

-- 3. Ноги — 25.04
WITH w AS (
  INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Ноги', '2026-04-25', 60, 'Выпады даются тяжело'
  FROM users WHERE login = 'elena_kozlova' RETURNING id
)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.exercise_id, e.sets, e.reps, e.weight_kg
FROM w, (VALUES (33,3,12,6.0),(35,3,15,20.0),(34,3,15,20.0),(36,3,15,0.0),(37,4,25,0.0)) AS e(exercise_id,sets,reps,weight_kg);

-- 4. Фулл-боди — 05.05
WITH w AS (
  INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Фулл-боди', '2026-05-05', 65, 'Чувствую прогресс!'
  FROM users WHERE login = 'elena_kozlova' RETURNING id
)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.exercise_id, e.sets, e.reps, e.weight_kg
FROM w, (VALUES (5,3,15,0.0),(33,3,12,8.0),(2,3,12,10.0),(17,3,15,6.0),(39,3,20,0.0)) AS e(exercise_id,sets,reps,weight_kg);

-- 5. Верх тела — 16.05
WITH w AS (
  INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Верх тела', '2026-05-16', 55, ''
  FROM users WHERE login = 'elena_kozlova' RETURNING id
)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.exercise_id, e.sets, e.reps, e.weight_kg
FROM w, (VALUES (5,3,15,0.0),(2,3,12,10.0),(10,3,12,10.0),(22,3,12,6.0),(40,3,45,0.0)) AS e(exercise_id,sets,reps,weight_kg);

-- 6. Пресс + Кардио — 28.05
WITH w AS (
  INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Пресс + Кардио', '2026-05-28', 45, 'Планка 60 сек — рекорд!'
  FROM users WHERE login = 'elena_kozlova' RETURNING id
)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.exercise_id, e.sets, e.reps, e.weight_kg
FROM w, (VALUES (49,1,30,0.0),(39,4,20,0.0),(40,3,60,0.0),(41,3,15,0.0),(43,3,30,0.0)) AS e(exercise_id,sets,reps,weight_kg);

-- ── Trainer comments ──────────────────────────────────────────────────────────
-- Comment on alexei's first chest workout
UPDATE workouts SET trainer_comment = 'Отличный результат! На следующей тренировке попробуй добавить 2.5 кг на жим лёжа. Следи за траекторией локтей.'
WHERE user_id = (SELECT id FROM users WHERE login = 'alexei_ivanov')
  AND date = '2026-04-07' AND title = 'Грудь + Трицепс';

-- Comment on alexei's second chest workout
UPDATE workouts SET trainer_comment = 'Прогресс виден! Продолжай в том же темпе. Добавь растяжку грудных после тренировки.'
WHERE user_id = (SELECT id FROM users WHERE login = 'alexei_ivanov')
  AND date = '2026-05-05' AND title = 'Грудь + Трицепс';

-- Comment on dmitry's heavy squat PR
UPDATE workouts SET trainer_comment = 'Великолепный присед! Техника чистая. Следи за нейтральным положением спины при 135+ кг.'
WHERE user_id = (SELECT id FROM users WHERE login = 'dmitry_smirnov')
  AND date = '2026-06-02' AND title = 'Ноги — тяжёлый день';

-- Comment on maria's full-body workout
UPDATE workouts SET trainer_comment = 'Хорошая работа! Рекомендую добавить 1 подход на бицепс — пока не хватает объёма для рук.'
WHERE user_id = (SELECT id FROM users WHERE login = 'maria_petrova')
  AND date = '2026-04-30' AND title = 'Верх тела';

-- ── Schedule entries ──────────────────────────────────────────────────────────
-- completed: dmitry — 28.05
INSERT INTO schedule (trainer_id, client_id, title, scheduled_at, duration_minutes, status, notes)
SELECT
  (SELECT id FROM users WHERE login = 'macaretchurilov'),
  (SELECT id FROM users WHERE login = 'dmitry_smirnov'),
  'Ноги + разбор техники',
  '2026-05-28 10:00:00+03',
  90,
  'completed',
  'Разобрали технику становой тяги';

-- planned: alexei — 10.06
INSERT INTO schedule (trainer_id, client_id, title, scheduled_at, duration_minutes, status, notes)
SELECT
  (SELECT id FROM users WHERE login = 'macaretchurilov'),
  (SELECT id FROM users WHERE login = 'alexei_ivanov'),
  'Грудь + Трицепс',
  '2026-06-10 11:00:00+03',
  75,
  'planned',
  'Контрольный замер максимума жима';

-- planned: maria — 12.06
INSERT INTO schedule (trainer_id, client_id, title, scheduled_at, duration_minutes, status, notes)
SELECT
  (SELECT id FROM users WHERE login = 'macaretchurilov'),
  (SELECT id FROM users WHERE login = 'maria_petrova'),
  'Фулл-боди + кардио',
  '2026-06-12 09:30:00+03',
  60,
  'planned',
  '';

-- ── Workout templates ─────────────────────────────────────────────────────────
-- Template 1: alexei — powerlifting upper body
WITH t AS (
  INSERT INTO workout_templates (user_id, name, is_public)
  SELECT id, 'Грудь + Трицепс (базовый)', true
  FROM users WHERE login = 'alexei_ivanov' RETURNING id
)
INSERT INTO template_exercises (template_id, exercise_id, sets, reps, weight_kg)
SELECT t.id, e.exercise_id, e.sets, e.reps, e.weight_kg
FROM t, (VALUES (1,4,8,100.0),(2,4,8,32.0),(26,4,8,70.0),(28,3,12,27.5),(29,3,10,0.0)) AS e(exercise_id,sets,reps,weight_kg);

-- Template 2: alexei — leg day
WITH t AS (
  INSERT INTO workout_templates (user_id, name, is_public)
  SELECT id, 'День ног (базовый)', true
  FROM users WHERE login = 'alexei_ivanov' RETURNING id
)
INSERT INTO template_exercises (template_id, exercise_id, sets, reps, weight_kg)
SELECT t.id, e.exercise_id, e.sets, e.reps, e.weight_kg
FROM t, (VALUES (31,5,5,110.0),(32,4,12,150.0),(14,3,5,120.0),(34,3,15,30.0),(37,4,20,0.0)) AS e(exercise_id,sets,reps,weight_kg);

-- Template 3: maria — full-body
WITH t AS (
  INSERT INTO workout_templates (user_id, name, is_public)
  SELECT id, 'Фулл-боди (женский)', false
  FROM users WHERE login = 'maria_petrova' RETURNING id
)
INSERT INTO template_exercises (template_id, exercise_id, sets, reps, weight_kg)
SELECT t.id, e.exercise_id, e.sets, e.reps, e.weight_kg
FROM t, (VALUES (5,3,15,0.0),(33,3,12,10.0),(2,3,12,16.0),(17,3,15,8.0),(39,3,20,0.0)) AS e(exercise_id,sets,reps,weight_kg);

-- Template 4: dmitry — heavy back
WITH t AS (
  INSERT INTO workout_templates (user_id, name, is_public)
  SELECT id, 'Спина — тяжёлый день', true
  FROM users WHERE login = 'dmitry_smirnov' RETURNING id
)
INSERT INTO template_exercises (template_id, exercise_id, sets, reps, weight_kg)
SELECT t.id, e.exercise_id, e.sets, e.reps, e.weight_kg
FROM t, (VALUES (14,5,5,160.0),(9,4,8,110.0),(8,3,8,0.0),(10,3,10,40.0),(13,3,15,50.0)) AS e(exercise_id,sets,reps,weight_kg);

COMMIT;
