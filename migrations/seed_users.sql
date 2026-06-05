-- Seed: coach (admin) + user1..user5, workouts, body_metrics, templates, schedule.
-- Password for ALL users: 1234 (bcrypt cost=10)

BEGIN;

-- ══════════════════════════════════════════════════════════════════
-- 1. ПОЛЬЗОВАТЕЛИ
-- ══════════════════════════════════════════════════════════════════
INSERT INTO users (login, email, password_hash, role) VALUES
  ('coach', 'coach@gym.local',  '$2a$10$4PjgFy0ScH8LuwgvHvYkBu2wxBo7hPg8VrSluXYwDaam2xzqrFOgu', 'admin'),
  ('user1', 'user1@gym.local',  '$2a$10$4PjgFy0ScH8LuwgvHvYkBu2wxBo7hPg8VrSluXYwDaam2xzqrFOgu', 'user'),
  ('user2', 'user2@gym.local',  '$2a$10$4PjgFy0ScH8LuwgvHvYkBu2wxBo7hPg8VrSluXYwDaam2xzqrFOgu', 'user'),
  ('user3', 'user3@gym.local',  '$2a$10$4PjgFy0ScH8LuwgvHvYkBu2wxBo7hPg8VrSluXYwDaam2xzqrFOgu', 'user'),
  ('user4', 'user4@gym.local',  '$2a$10$4PjgFy0ScH8LuwgvHvYkBu2wxBo7hPg8VrSluXYwDaam2xzqrFOgu', 'user'),
  ('user5', 'user5@gym.local',  '$2a$10$4PjgFy0ScH8LuwgvHvYkBu2wxBo7hPg8VrSluXYwDaam2xzqrFOgu', 'user')
ON CONFLICT (login) DO NOTHING;

-- ══════════════════════════════════════════════════════════════════
-- 2. ТРЕНИРОВКИ — user1 (силовик, мужчина ~85 кг)
-- ══════════════════════════════════════════════════════════════════
WITH w AS (INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Грудь + Трицепс', '2026-04-01', 75, 'Жим идёт хорошо'
  FROM users WHERE login = 'user1' RETURNING id)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.i, e.s, e.r, e.w FROM w,
  (VALUES (1,4,8,100.0),(3,4,8,80.0),(2,3,10,32.0),(25,3,10,70.0),(27,4,12,32.5)) AS e(i,s,r,w);

WITH w AS (INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Спина + Бицепс', '2026-04-03', 80, ''
  FROM users WHERE login = 'user1' RETURNING id)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.i, e.s, e.r, e.w FROM w,
  (VALUES (9,4,8,90.0),(11,3,12,55.0),(10,3,10,28.0),(20,4,10,45.0),(22,3,12,16.0)) AS e(i,s,r,w);

WITH w AS (INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'День ног', '2026-04-07', 90, 'Присед 110 кг × 5'
  FROM users WHERE login = 'user1' RETURNING id)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.i, e.s, e.r, e.w FROM w,
  (VALUES (30,5,5,110.0),(31,4,12,150.0),(37,3,5,120.0),(33,3,15,30.0),(36,4,20,0.0)) AS e(i,s,r,w);

WITH w AS (INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Плечи', '2026-04-10', 65, ''
  FROM users WHERE login = 'user1' RETURNING id)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.i, e.s, e.r, e.w FROM w,
  (VALUES (14,4,8,60.0),(15,3,10,22.0),(16,4,15,12.0),(17,3,12,40.0),(19,3,10,20.0)) AS e(i,s,r,w);

WITH w AS (INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Грудь + Трицепс', '2026-04-15', 80, 'Жим 102.5 кг — личник!'
  FROM users WHERE login = 'user1' RETURNING id)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.i, e.s, e.r, e.w FROM w,
  (VALUES (1,4,8,102.5),(3,4,8,82.5),(2,3,10,34.0),(26,4,8,40.0),(27,4,12,35.0)) AS e(i,s,r,w);

WITH w AS (INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Спина + Бицепс', '2026-04-17', 75, 'Тяга — добавил вес'
  FROM users WHERE login = 'user1' RETURNING id)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.i, e.s, e.r, e.w FROM w,
  (VALUES (9,4,8,92.5),(11,3,12,57.5),(12,3,12,50.0),(20,4,10,47.5),(21,3,12,18.0)) AS e(i,s,r,w);

WITH w AS (INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'День ног', '2026-04-22', 95, 'Становая 125 кг — чисто!'
  FROM users WHERE login = 'user1' RETURNING id)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.i, e.s, e.r, e.w FROM w,
  (VALUES (30,5,5,112.5),(37,4,5,125.0),(31,4,12,155.0),(33,3,15,32.5),(34,3,15,35.0)) AS e(i,s,r,w);

WITH w AS (INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Грудь + Трицепс', '2026-04-28', 80, ''
  FROM users WHERE login = 'user1' RETURNING id)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.i, e.s, e.r, e.w FROM w,
  (VALUES (1,4,8,105.0),(3,4,8,85.0),(2,3,10,36.0),(25,3,8,75.0),(27,4,12,37.5)) AS e(i,s,r,w);

-- ══════════════════════════════════════════════════════════════════
-- 3. ТРЕНИРОВКИ — user2 (фитнес / похудение, женщина ~62 кг)
-- ══════════════════════════════════════════════════════════════════
WITH w AS (INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Фулл-боди', '2026-04-05', 60, 'Первая тренировка'
  FROM users WHERE login = 'user2' RETURNING id)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.i, e.s, e.r, e.w FROM w,
  (VALUES (5,3,15,0.0),(32,3,12,8.0),(16,3,15,6.0),(39,3,20,0.0),(40,1,30,0.0)) AS e(i,s,r,w);

WITH w AS (INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Кардио + Пресс', '2026-04-08', 45, 'Кардио 20 мин + пресс'
  FROM users WHERE login = 'user2' RETURNING id)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.i, e.s, e.r, e.w FROM w,
  (VALUES (45,1,20,0.0),(49,1,15,0.0),(39,4,20,0.0),(40,3,60,0.0),(41,3,15,0.0)) AS e(i,s,r,w);

WITH w AS (INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Фулл-боди', '2026-04-12', 65, 'Добавила вес на выпады'
  FROM users WHERE login = 'user2' RETURNING id)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.i, e.s, e.r, e.w FROM w,
  (VALUES (5,3,15,0.0),(32,3,12,10.0),(2,3,12,14.0),(16,3,15,7.0),(39,3,20,0.0)) AS e(i,s,r,w);

WITH w AS (INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Верх тела', '2026-04-16', 55, 'Прогресс на жиме!'
  FROM users WHERE login = 'user2' RETURNING id)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.i, e.s, e.r, e.w FROM w,
  (VALUES (2,3,12,16.0),(10,3,12,14.0),(16,4,15,8.0),(21,3,12,8.0),(40,4,45,0.0)) AS e(i,s,r,w);

WITH w AS (INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Кардио + Пресс', '2026-04-21', 50, ''
  FROM users WHERE login = 'user2' RETURNING id)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.i, e.s, e.r, e.w FROM w,
  (VALUES (47,1,10,0.0),(45,1,25,0.0),(39,4,20,0.0),(41,3,15,0.0),(43,3,20,0.0)) AS e(i,s,r,w);

WITH w AS (INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Фулл-боди', '2026-04-28', 60, 'Стабильно, без рекордов'
  FROM users WHERE login = 'user2' RETURNING id)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.i, e.s, e.r, e.w FROM w,
  (VALUES (5,3,15,0.0),(32,3,12,10.0),(2,3,12,16.0),(39,4,20,0.0),(40,3,60,0.0)) AS e(i,s,r,w);

-- ══════════════════════════════════════════════════════════════════
-- 4. ТРЕНИРОВКИ — user3 (пауэрлифтер, мужчина ~95 кг)
-- ══════════════════════════════════════════════════════════════════
WITH w AS (INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Ноги — тяжёлый день', '2026-04-06', 95, 'Присед 130 кг × 5'
  FROM users WHERE login = 'user3' RETURNING id)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.i, e.s, e.r, e.w FROM w,
  (VALUES (30,5,5,130.0),(37,3,5,150.0),(31,4,10,180.0),(33,3,12,45.0),(36,5,20,0.0)) AS e(i,s,r,w);

WITH w AS (INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Грудь — тяжёлый день', '2026-04-09', 85, 'Жим 125 кг — новый максимум!'
  FROM users WHERE login = 'user3' RETURNING id)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.i, e.s, e.r, e.w FROM w,
  (VALUES (1,5,5,125.0),(3,3,8,95.0),(2,3,10,40.0),(28,3,10,0.0),(25,3,8,80.0)) AS e(i,s,r,w);

WITH w AS (INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Спина — тяжёлый день', '2026-04-13', 90, 'Становая 160 кг!'
  FROM users WHERE login = 'user3' RETURNING id)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.i, e.s, e.r, e.w FROM w,
  (VALUES (37,5,5,160.0),(9,4,8,110.0),(8,3,8,0.0),(10,3,10,40.0),(13,3,15,50.0)) AS e(i,s,r,w);

WITH w AS (INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Ноги — тяжёлый день', '2026-04-20', 90, ''
  FROM users WHERE login = 'user3' RETURNING id)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.i, e.s, e.r, e.w FROM w,
  (VALUES (30,5,5,132.5),(37,3,5,155.0),(31,4,10,185.0),(38,3,10,80.0),(36,5,20,0.0)) AS e(i,s,r,w);

WITH w AS (INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Грудь + Плечи', '2026-04-25', 80, 'Армейский жим — прогресс'
  FROM users WHERE login = 'user3' RETURNING id)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.i, e.s, e.r, e.w FROM w,
  (VALUES (1,4,8,115.0),(14,4,5,80.0),(3,3,10,90.0),(16,4,12,20.0),(17,3,10,55.0)) AS e(i,s,r,w);

WITH w AS (INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Спина', '2026-04-30', 75, ''
  FROM users WHERE login = 'user3' RETURNING id)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.i, e.s, e.r, e.w FROM w,
  (VALUES (37,3,5,160.0),(9,4,8,112.5),(11,3,12,62.5),(12,3,12,55.0),(13,3,15,55.0)) AS e(i,s,r,w);

WITH w AS (INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Ноги — тяжёлый день', '2026-05-05', 95, 'Присед 135 кг — личник!'
  FROM users WHERE login = 'user3' RETURNING id)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.i, e.s, e.r, e.w FROM w,
  (VALUES (30,5,5,135.0),(37,4,5,162.5),(31,4,10,190.0),(33,3,12,47.5),(36,5,20,0.0)) AS e(i,s,r,w);

-- ══════════════════════════════════════════════════════════════════
-- 5. ТРЕНИРОВКИ — user4 (лёгкий фитнес, женщина ~56 кг)
-- ══════════════════════════════════════════════════════════════════
WITH w AS (INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Верх тела', '2026-04-07', 50, 'Лёгкие веса, работаю над техникой'
  FROM users WHERE login = 'user4' RETURNING id)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.i, e.s, e.r, e.w FROM w,
  (VALUES (5,3,12,0.0),(2,3,12,8.0),(16,3,15,5.0),(21,3,12,5.0),(39,3,30,0.0)) AS e(i,s,r,w);

WITH w AS (INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Пресс + Кардио', '2026-04-10', 40, ''
  FROM users WHERE login = 'user4' RETURNING id)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.i, e.s, e.r, e.w FROM w,
  (VALUES (49,1,25,0.0),(39,3,20,0.0),(40,3,30,0.0),(41,3,12,0.0),(42,3,20,0.0)) AS e(i,s,r,w);

WITH w AS (INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Ноги', '2026-04-14', 55, 'Выпады с гантелями'
  FROM users WHERE login = 'user4' RETURNING id)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.i, e.s, e.r, e.w FROM w,
  (VALUES (32,3,12,6.0),(34,3,15,20.0),(33,3,15,20.0),(35,3,15,0.0),(36,4,25,0.0)) AS e(i,s,r,w);

WITH w AS (INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Фулл-боди', '2026-04-19', 60, 'Чувствую прогресс!'
  FROM users WHERE login = 'user4' RETURNING id)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.i, e.s, e.r, e.w FROM w,
  (VALUES (5,3,15,0.0),(32,3,12,8.0),(2,3,12,10.0),(16,3,15,6.0),(39,3,20,0.0)) AS e(i,s,r,w);

WITH w AS (INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Верх тела', '2026-04-24', 55, ''
  FROM users WHERE login = 'user4' RETURNING id)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.i, e.s, e.r, e.w FROM w,
  (VALUES (5,3,15,0.0),(2,3,12,10.0),(10,3,12,10.0),(21,3,12,6.0),(40,3,45,0.0)) AS e(i,s,r,w);

WITH w AS (INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Пресс + Кардио', '2026-05-01', 45, 'Планка 60 сек — рекорд!'
  FROM users WHERE login = 'user4' RETURNING id)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.i, e.s, e.r, e.w FROM w,
  (VALUES (49,1,30,0.0),(39,4,20,0.0),(40,3,60,0.0),(41,3,15,0.0),(43,3,30,0.0)) AS e(i,s,r,w);

-- ══════════════════════════════════════════════════════════════════
-- 6. ТРЕНИРОВКИ — user5 (новичок, мужчина ~75 кг)
-- ══════════════════════════════════════════════════════════════════
WITH w AS (INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Первая тренировка', '2026-04-10', 45, 'Начинаю заниматься'
  FROM users WHERE login = 'user5' RETURNING id)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.i, e.s, e.r, e.w FROM w,
  (VALUES (5,3,10,0.0),(30,3,10,60.0),(8,3,8,0.0),(39,3,15,0.0),(40,1,20,0.0)) AS e(i,s,r,w);

WITH w AS (INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Фулл-боди', '2026-04-14', 50, 'Уже легче!'
  FROM users WHERE login = 'user5' RETURNING id)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.i, e.s, e.r, e.w FROM w,
  (VALUES (5,3,12,0.0),(30,3,10,65.0),(9,3,8,50.0),(20,3,10,30.0),(39,3,20,0.0)) AS e(i,s,r,w);

WITH w AS (INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Грудь + Плечи', '2026-04-18', 55, ''
  FROM users WHERE login = 'user5' RETURNING id)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.i, e.s, e.r, e.w FROM w,
  (VALUES (1,3,8,60.0),(15,3,10,14.0),(5,3,12,0.0),(16,3,12,8.0),(40,2,30,0.0)) AS e(i,s,r,w);

WITH w AS (INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Спина + Бицепс', '2026-04-23', 60, 'Первый раз подтянулся!'
  FROM users WHERE login = 'user5' RETURNING id)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.i, e.s, e.r, e.w FROM w,
  (VALUES (8,3,5,0.0),(11,3,10,40.0),(10,3,10,14.0),(20,3,10,30.0),(22,3,12,12.0)) AS e(i,s,r,w);

WITH w AS (INSERT INTO workouts (user_id, title, date, duration_minutes, notes)
  SELECT id, 'Ноги', '2026-04-28', 65, 'Присед 70 кг!'
  FROM users WHERE login = 'user5' RETURNING id)
INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight_kg)
SELECT w.id, e.i, e.s, e.r, e.w FROM w,
  (VALUES (30,3,8,70.0),(31,3,10,80.0),(32,3,10,10.0),(33,3,12,25.0),(36,3,20,0.0)) AS e(i,s,r,w);

-- ══════════════════════════════════════════════════════════════════
-- 7. КОММЕНТАРИИ ТРЕНЕРА
-- ══════════════════════════════════════════════════════════════════
UPDATE workouts SET trainer_comment = 'Отличный прогресс! На следующей тренировке добавь 2.5 кг. Следи за траекторией локтей — разводи чуть шире.'
WHERE user_id = (SELECT id FROM users WHERE login = 'user1')
  AND date = '2026-04-15' AND title = 'Грудь + Трицепс';

UPDATE workouts SET trainer_comment = 'Становая сделана чисто! Продолжай работу над техникой. При 150+ кг обязательно используй пояс.'
WHERE user_id = (SELECT id FROM users WHERE login = 'user3')
  AND date = '2026-05-05' AND title = 'Ноги — тяжёлый день';

UPDATE workouts SET trainer_comment = 'Хорошая работа! Рекомендую добавить 1 подход на трицепс — пока не хватает объёма для рук.'
WHERE user_id = (SELECT id FROM users WHERE login = 'user2')
  AND date = '2026-04-16' AND title = 'Верх тела';

UPDATE workouts SET trainer_comment = 'Молодец! Планка 60 сек — это уже хороший результат. Следующая цель — 90 сек.'
WHERE user_id = (SELECT id FROM users WHERE login = 'user4')
  AND date = '2026-05-01' AND title = 'Пресс + Кардио';

UPDATE workouts SET trainer_comment = 'Отличное начало! Не торопись с весами — сначала техника. На следующей сессии разберём присед подробнее.'
WHERE user_id = (SELECT id FROM users WHERE login = 'user5')
  AND date = '2026-04-10' AND title = 'Первая тренировка';

-- ══════════════════════════════════════════════════════════════════
-- 8. ЗАМЕРЫ ТЕЛА (body_metrics)
-- ══════════════════════════════════════════════════════════════════
-- user1 (мужчина, силовик)
INSERT INTO body_metrics (user_id, weight_kg, body_fat_percent, chest_cm, waist_cm, hips_cm, bicep_cm, measured_at)
SELECT id, 85.5, 15.2, 105.0, 84.0, 97.0, 38.0, '2026-04-01' FROM users WHERE login = 'user1';
INSERT INTO body_metrics (user_id, weight_kg, body_fat_percent, chest_cm, waist_cm, hips_cm, bicep_cm, measured_at)
SELECT id, 86.0, 14.8, 106.0, 83.5, 97.5, 38.5, '2026-04-15' FROM users WHERE login = 'user1';
INSERT INTO body_metrics (user_id, weight_kg, body_fat_percent, chest_cm, waist_cm, hips_cm, bicep_cm, measured_at)
SELECT id, 86.5, 14.5, 107.0, 83.0, 98.0, 39.0, '2026-05-01' FROM users WHERE login = 'user1';

-- user2 (женщина, похудение)
INSERT INTO body_metrics (user_id, weight_kg, body_fat_percent, chest_cm, waist_cm, hips_cm, bicep_cm, measured_at)
SELECT id, 65.0, 27.5, 90.0, 72.0, 98.0, 27.0, '2026-04-05' FROM users WHERE login = 'user2';
INSERT INTO body_metrics (user_id, weight_kg, body_fat_percent, chest_cm, waist_cm, hips_cm, bicep_cm, measured_at)
SELECT id, 64.2, 26.8, 89.5, 71.0, 97.0, 26.5, '2026-04-19' FROM users WHERE login = 'user2';
INSERT INTO body_metrics (user_id, weight_kg, body_fat_percent, chest_cm, waist_cm, hips_cm, bicep_cm, measured_at)
SELECT id, 63.5, 26.0, 89.0, 70.0, 96.0, 26.5, '2026-05-03' FROM users WHERE login = 'user2';
INSERT INTO body_metrics (user_id, weight_kg, body_fat_percent, chest_cm, waist_cm, hips_cm, bicep_cm, measured_at)
SELECT id, 62.8, 25.2, 88.5, 69.5, 95.5, 26.0, '2026-05-17' FROM users WHERE login = 'user2';

-- user3 (мужчина, пауэрлифтер)
INSERT INTO body_metrics (user_id, weight_kg, body_fat_percent, chest_cm, waist_cm, hips_cm, bicep_cm, measured_at)
SELECT id, 95.0, 18.0, 115.0, 92.0, 106.0, 42.0, '2026-04-01' FROM users WHERE login = 'user3';
INSERT INTO body_metrics (user_id, weight_kg, body_fat_percent, chest_cm, waist_cm, hips_cm, bicep_cm, measured_at)
SELECT id, 95.5, 17.8, 115.5, 91.5, 106.5, 42.5, '2026-04-20' FROM users WHERE login = 'user3';
INSERT INTO body_metrics (user_id, weight_kg, body_fat_percent, chest_cm, waist_cm, hips_cm, bicep_cm, measured_at)
SELECT id, 96.0, 17.5, 116.0, 91.0, 107.0, 43.0, '2026-05-05' FROM users WHERE login = 'user3';

-- user4 (женщина, фитнес)
INSERT INTO body_metrics (user_id, weight_kg, body_fat_percent, chest_cm, waist_cm, hips_cm, bicep_cm, measured_at)
SELECT id, 58.0, 24.0, 87.0, 68.0, 94.0, 25.5, '2026-04-07' FROM users WHERE login = 'user4';
INSERT INTO body_metrics (user_id, weight_kg, body_fat_percent, chest_cm, waist_cm, hips_cm, bicep_cm, measured_at)
SELECT id, 57.5, 23.5, 86.5, 67.5, 93.5, 25.5, '2026-04-21' FROM users WHERE login = 'user4';
INSERT INTO body_metrics (user_id, weight_kg, body_fat_percent, chest_cm, waist_cm, hips_cm, bicep_cm, measured_at)
SELECT id, 57.0, 23.0, 86.0, 67.0, 93.0, 25.0, '2026-05-05' FROM users WHERE login = 'user4';

-- user5 (мужчина, новичок)
INSERT INTO body_metrics (user_id, weight_kg, body_fat_percent, chest_cm, waist_cm, hips_cm, bicep_cm, measured_at)
SELECT id, 75.0, 22.0, 96.0, 82.0, 98.0, 33.0, '2026-04-10' FROM users WHERE login = 'user5';
INSERT INTO body_metrics (user_id, weight_kg, body_fat_percent, chest_cm, waist_cm, hips_cm, bicep_cm, measured_at)
SELECT id, 75.5, 21.5, 96.5, 81.5, 98.0, 33.5, '2026-04-28' FROM users WHERE login = 'user5';

-- ══════════════════════════════════════════════════════════════════
-- 9. ШАБЛОНЫ ТРЕНИРОВОК
-- ══════════════════════════════════════════════════════════════════
WITH t AS (INSERT INTO workout_templates (user_id, name, is_public)
  SELECT id, 'Грудь + Трицепс (базовый)', true FROM users WHERE login = 'user1' RETURNING id)
INSERT INTO template_exercises (template_id, exercise_id, sets, reps, weight_kg)
SELECT t.id, e.i, e.s, e.r, e.w FROM t,
  (VALUES (1,4,8,100.0),(3,4,8,80.0),(25,3,8,70.0),(27,4,12,32.5),(28,3,10,0.0)) AS e(i,s,r,w);

WITH t AS (INSERT INTO workout_templates (user_id, name, is_public)
  SELECT id, 'День ног', true FROM users WHERE login = 'user1' RETURNING id)
INSERT INTO template_exercises (template_id, exercise_id, sets, reps, weight_kg)
SELECT t.id, e.i, e.s, e.r, e.w FROM t,
  (VALUES (30,5,5,110.0),(37,4,5,120.0),(31,4,12,150.0),(33,3,15,30.0),(36,4,20,0.0)) AS e(i,s,r,w);

WITH t AS (INSERT INTO workout_templates (user_id, name, is_public)
  SELECT id, 'Фулл-боди (женский)', false FROM users WHERE login = 'user2' RETURNING id)
INSERT INTO template_exercises (template_id, exercise_id, sets, reps, weight_kg)
SELECT t.id, e.i, e.s, e.r, e.w FROM t,
  (VALUES (5,3,15,0.0),(32,3,12,10.0),(2,3,12,16.0),(16,3,15,8.0),(39,3,20,0.0)) AS e(i,s,r,w);

WITH t AS (INSERT INTO workout_templates (user_id, name, is_public)
  SELECT id, 'Спина — тяжёлый день', true FROM users WHERE login = 'user3' RETURNING id)
INSERT INTO template_exercises (template_id, exercise_id, sets, reps, weight_kg)
SELECT t.id, e.i, e.s, e.r, e.w FROM t,
  (VALUES (37,5,5,160.0),(9,4,8,110.0),(8,3,8,0.0),(10,3,10,40.0),(13,3,15,50.0)) AS e(i,s,r,w);

WITH t AS (INSERT INTO workout_templates (user_id, name, is_public)
  SELECT id, 'Фулл-боди новичок', true FROM users WHERE login = 'user5' RETURNING id)
INSERT INTO template_exercises (template_id, exercise_id, sets, reps, weight_kg)
SELECT t.id, e.i, e.s, e.r, e.w FROM t,
  (VALUES (5,3,12,0.0),(30,3,10,60.0),(8,3,8,0.0),(20,3,10,25.0),(39,3,15,0.0)) AS e(i,s,r,w);

-- ══════════════════════════════════════════════════════════════════
-- 10. РАСПИСАНИЕ ТРЕНЕРА (coach)
-- ══════════════════════════════════════════════════════════════════
-- завершённые сессии
INSERT INTO schedule (trainer_id, client_id, title, scheduled_at, duration_minutes, status, notes) SELECT
  (SELECT id FROM users WHERE login = 'coach'),
  (SELECT id FROM users WHERE login = 'user5'),
  'Вводная тренировка',
  '2026-04-09 12:00:00+03', 60, 'completed', 'Знакомство с оборудованием. Базовые движения освоены';

INSERT INTO schedule (trainer_id, client_id, title, scheduled_at, duration_minutes, status, notes) SELECT
  (SELECT id FROM users WHERE login = 'coach'),
  (SELECT id FROM users WHERE login = 'user2'),
  'Фулл-боди + замер прогресса',
  '2026-05-17 09:30:00+03', 60, 'completed', 'Отличный прогресс! −2 кг за месяц';

INSERT INTO schedule (trainer_id, client_id, title, scheduled_at, duration_minutes, status, notes) SELECT
  (SELECT id FROM users WHERE login = 'coach'),
  (SELECT id FROM users WHERE login = 'user3'),
  'Ноги + разбор техники становой',
  '2026-05-20 10:00:00+03', 90, 'completed', 'Разобрали технику становой тяги. Ставить пояс при 150+ кг';

INSERT INTO schedule (trainer_id, client_id, title, scheduled_at, duration_minutes, status, notes) SELECT
  (SELECT id FROM users WHERE login = 'coach'),
  (SELECT id FROM users WHERE login = 'user1'),
  'Грудь + контроль техники',
  '2026-05-22 11:00:00+03', 75, 'completed', 'Проверили технику жима. Рекомендовано больше разминки';

-- отменённая сессия
INSERT INTO schedule (trainer_id, client_id, title, scheduled_at, duration_minutes, status, notes) SELECT
  (SELECT id FROM users WHERE login = 'coach'),
  (SELECT id FROM users WHERE login = 'user4'),
  'Ноги',
  '2026-05-28 11:00:00+03', 55, 'cancelled', 'Клиент заболел';

-- запланированные сессии
INSERT INTO schedule (trainer_id, client_id, title, scheduled_at, duration_minutes, status, notes) SELECT
  (SELECT id FROM users WHERE login = 'coach'),
  (SELECT id FROM users WHERE login = 'user1'),
  'Грудь + Трицепс',
  '2026-06-10 11:00:00+03', 75, 'planned', 'Контрольный замер максимума жима';

INSERT INTO schedule (trainer_id, client_id, title, scheduled_at, duration_minutes, status, notes) SELECT
  (SELECT id FROM users WHERE login = 'coach'),
  (SELECT id FROM users WHERE login = 'user2'),
  'Фулл-боди + кардио',
  '2026-06-12 09:30:00+03', 60, 'planned', '';

INSERT INTO schedule (trainer_id, client_id, title, scheduled_at, duration_minutes, status, notes) SELECT
  (SELECT id FROM users WHERE login = 'coach'),
  (SELECT id FROM users WHERE login = 'user3'),
  'Ноги — тяжёлый день',
  '2026-06-13 10:00:00+03', 90, 'planned', 'Работа над приседом 140 кг';

INSERT INTO schedule (trainer_id, client_id, title, scheduled_at, duration_minutes, status, notes) SELECT
  (SELECT id FROM users WHERE login = 'coach'),
  (SELECT id FROM users WHERE login = 'user4'),
  'Верх тела + кор',
  '2026-06-14 10:00:00+03', 55, 'planned', 'Добавим упражнения на спину';

INSERT INTO schedule (trainer_id, client_id, title, scheduled_at, duration_minutes, status, notes) SELECT
  (SELECT id FROM users WHERE login = 'coach'),
  (SELECT id FROM users WHERE login = 'user5'),
  'Прогресс-тренировка',
  '2026-06-16 16:00:00+03', 60, 'planned', 'Проверим прогресс за месяц';

COMMIT;
