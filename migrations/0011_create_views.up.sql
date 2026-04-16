BEGIN;

-- CREATE BASE VIEW

CREATE OR REPLACE VIEW vw_schedule_detailed_base AS
SELECT
    s.id AS schedule_id,
    s.semester_id,
    s.version,
    s.school_id,

    r.id AS requirement_id,
    r.teacher_id,
    r.subject_id,
    r.class_id,

    t.last_name || ' ' || t.first_name AS teacher_name,
    sub.name AS subject_name,
    c.grade,
    c.class AS class_name,

    rm.id AS room_id,
    rm.name AS room_name,

    ts.day_of_week,
    ts.start_time,
    ts.end_time

FROM schedules s
JOIN requirements r   ON s.requirement_id = r.id
JOIN teachers t       ON r.teacher_id = t.id
JOIN subjects sub     ON r.subject_id = sub.id
JOIN classes c        ON r.class_id = c.id
LEFT JOIN rooms rm    ON s.room_id = rm.id
JOIN timeslots ts     ON s.timeslot_id = ts.id;

-- Weekly classes schedule view

CREATE OR REPLACE VIEW vw_class_weekly_schedule AS
SELECT
    semester_id,
    version,
    class_id,
    grade,
    class_name,

    day_of_week,
    start_time,
    end_time,

    subject_name,
    teacher_name,
    room_name

FROM vw_schedule_detailed_base;

-- Teacher weekly schedule view

CREATE OR REPLACE VIEW vw_teacher_weekly_schedule AS
SELECT
    semester_id,
    version,
    teacher_id,
    teacher_name,

    day_of_week,
    start_time,
    end_time,

    subject_name,
    grade,
    class_name,
    room_name

FROM vw_schedule_detailed_base;

-- Room utilization view

CREATE OR REPLACE VIEW vw_room_schedule AS
SELECT
    semester_id,
    version,
    room_id,
    room_name,

    day_of_week,
    start_time,
    end_time,

    subject_name,
    teacher_name,
    grade,
    class_name

FROM vw_schedule_detailed_base
WHERE room_id IS NOT NULL;

-- Teacher confliction detection view

CREATE OR REPLACE VIEW vw_teacher_conflicts AS
SELECT
    teacher_id,
    day_of_week,
    start_time,
    COUNT(*) AS conflict_count
FROM vw_schedule_detailed_base
GROUP BY teacher_id, day_of_week, start_time
HAVING COUNT(*) > 1;

-- Room confliction detection view

CREATE OR REPLACE VIEW vw_room_conflicts AS
SELECT
    room_id,
    day_of_week,
    start_time,
    COUNT(*) AS conflict_count
FROM vw_schedule_detailed_base
WHERE room_id IS NOT NULL
GROUP BY room_id, day_of_week, start_time
HAVING COUNT(*) > 1;

-- Classes confliction detection view

CREATE OR REPLACE VIEW vw_class_conflicts AS
SELECT
    class_id,
    day_of_week,
    start_time,
    COUNT(*) AS conflict_count
FROM vw_schedule_detailed_base
GROUP BY class_id, day_of_week, start_time
HAVING COUNT(*) > 1;

-- Teacher weekly workload

CREATE OR REPLACE VIEW vw_teacher_load AS
SELECT
    semester_id,
    version,
    teacher_id,
    teacher_name,
    COUNT(*) AS total_sessions
FROM vw_schedule_detailed_base
GROUP BY semester_id, version, teacher_id, teacher_name;

-- Class subject distribution

CREATE OR REPLACE VIEW vw_class_subject_distribution AS
SELECT
    semester_id,
    version,
    class_id,
    subject_name,
    COUNT(*) AS sessions_per_week
FROM vw_schedule_detailed_base
GROUP BY semester_id, version, class_id, subject_name;


COMMIT;