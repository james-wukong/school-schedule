-- ============================================================================
-- SCHOOL SCHEDULING SYSTEM - PostgreSQL Schema
-- Version: 1.0
-- Tables: schools, teachers, students, courses, classrooms, timeslots, 
--         constraints, schedules, schedule_entries
-- ============================================================================

-- ============================================================================
-- SCHOOLS
-- ============================================================================

CREATE TABLE schools (
    id BIGINT GENERATED ALWAYS AS IDENTITY (START WITH 1000 INCREMENT BY 1) PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    code VARCHAR(50) NOT NULL UNIQUE,
    address TEXT,
    city VARCHAR(100),
    state VARCHAR(50),
    postal_code VARCHAR(20),
    country VARCHAR(100),
    phone VARCHAR(20),
    email VARCHAR(100),
    website VARCHAR(255),
    principal_name VARCHAR(255),
    established_year INT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_schools_code ON schools(code);
CREATE INDEX idx_schools_active ON schools(is_active);

-- ============================================================================
-- TEACHERS
-- ============================================================================

CREATE TABLE teachers (
    id BIGINT GENERATED ALWAYS AS IDENTITY (START WITH 1000 INCREMENT BY 1) PRIMARY KEY,
    school_id BIGINT NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    employee_id BIGINT NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(100),
    phone VARCHAR(20),
    qualification VARCHAR(255),
    specialization VARCHAR(100),
    hire_date DATE NOT NULL,
    employment_type VARCHAR(50), -- 'Full-time', 'Part-time', 'Contract'
    max_classes_per_day INT DEFAULT 5,
    available_days VARCHAR(50), -- 'Monday,Tuesday,Wednesday,Thursday,Friday'
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE(school_id, employee_id)
);

CREATE INDEX idx_teachers_school ON teachers(school_id);
CREATE INDEX idx_teachers_email ON teachers(email);
CREATE INDEX idx_teachers_active ON teachers(is_active);

-- ============================================================================
-- STUDENTS
-- ============================================================================

CREATE TABLE students (
    id BIGINT GENERATED ALWAYS AS IDENTITY (START WITH 1000 INCREMENT BY 1) PRIMARY KEY,
    school_id BIGINT NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    admission_number VARCHAR(50) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(100),
    phone VARCHAR(20),
    date_of_birth DATE,
    gender VARCHAR(20),
    blood_group VARCHAR(10),
    admission_date DATE NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE(school_id, admission_number)
);

CREATE INDEX idx_students_school ON students(school_id);
CREATE INDEX idx_students_admission ON students(admission_number);
CREATE INDEX idx_students_active ON students(is_active);

-- ============================================================================
-- COURSES
-- ============================================================================

CREATE TABLE courses (
    id BIGINT GENERATED ALWAYS AS IDENTITY (START WITH 1000 INCREMENT BY 1) PRIMARY KEY,
    school_id BIGINT NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    teacher_id BIGINT NOT NULL REFERENCES teachers(id) ON DELETE RESTRICT,
    semester INT, -- 1, 2, 3 etc
    credit_hours INT,
    max_students INT,
    duration_minutes INT DEFAULT 60, -- Class duration
    required_days_per_week INT DEFAULT 3, -- How many times per week
    min_day_gap INT DEFAULT 1, -- Minimum days between sessions
    preferred_days VARCHAR(100), -- 'Monday,Wednesday,Friday'
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE(school_id, code)
);

CREATE INDEX idx_courses_school ON courses(school_id);
CREATE INDEX idx_courses_teacher ON courses(teacher_id);
CREATE INDEX idx_courses_active ON courses(is_active);

-- ============================================================================
-- CLASSROOMS
-- ============================================================================

CREATE TABLE classrooms (
    id BIGINT GENERATED ALWAYS AS IDENTITY (START WITH 1000 INCREMENT BY 1) PRIMARY KEY,
    school_id BIGINT NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    room_type VARCHAR(50), -- 'Classroom', 'Lab', 'Auditorium', 'Library'
    capacity INT DEFAULT 40,
    floor_number INT,
    building VARCHAR(100),
    available_days VARCHAR(50), -- 'Monday,Tuesday,Wednesday,Thursday,Friday'
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE(school_id, code)
);

CREATE INDEX idx_classrooms_school ON classrooms(school_id);
CREATE INDEX idx_classrooms_type ON classrooms(room_type);
CREATE INDEX idx_classrooms_active ON classrooms(is_active);

-- ============================================================================
-- TIMESLOTS
-- ============================================================================

CREATE TABLE timeslots (
    id BIGINT GENERATED ALWAYS AS IDENTITY (START WITH 1000 INCREMENT BY 1) PRIMARY KEY,
    school_id BIGINT NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    day_of_week VARCHAR(20) NOT NULL, -- 'Monday', 'Tuesday', etc
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    slot_order INT, -- 1st slot, 2nd slot, etc
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(school_id, day_of_week, start_time)
);

CREATE INDEX idx_timeslots_school ON timeslots(school_id);
CREATE INDEX idx_timeslots_day ON timeslots(day_of_week);

-- ============================================================================
-- CONSTRAINTS
-- ============================================================================

CREATE TABLE constraints (
    id BIGINT GENERATED ALWAYS AS IDENTITY (START WITH 1000 INCREMENT BY 1) PRIMARY KEY,
    school_id BIGINT NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    constraint_type VARCHAR(50) NOT NULL, -- 'hard', 'soft'
    
    -- Constraint subject types
    subject_type VARCHAR(50), -- 'teacher', 'classroom', 'course', 'student'
    subject_id BIGINT, -- Reference to teacher, classroom, course, or student
    
    -- Constraint parameters
    description TEXT,
    rule_json JSONB, -- Flexible JSON for complex constraints
    weight INT DEFAULT 10, -- Weight for soft constraints
    
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_constraints_school ON constraints(school_id);
CREATE INDEX idx_constraints_type ON constraints(constraint_type);
CREATE INDEX idx_constraints_subject ON constraints(subject_type, subject_id);

-- Examples of constraint types (in rule_json):
-- {
--   type: no_consecutive_classes,
--   room_id: bigint,
--   min_break_minutes: 30
-- }
-- {
--   type: teacher_max_hours_per_day,
--   teacher_id: bigint,
--   max_hours: 6
-- }
-- {
--   type: preferred_timeslot,
--   timeslot_id: bigint,
--   priority: high
-- }
-- {
--   type: avoid_timeslot,
--   timeslot_id: bigint,
--   reason: teacher unavailable
-- }

-- ============================================================================
-- SCHEDULES
-- ============================================================================

CREATE TABLE schedules (
    id BIGINT GENERATED ALWAYS AS IDENTITY (START WITH 1000 INCREMENT BY 1) PRIMARY KEY,
    school_id BIGINT NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL, -- Fall 2024, Spring 2024
    code VARCHAR(50) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    status VARCHAR(50) DEFAULT 'Draft', -- 'Draft', 'Published', 'Active', 'Archived'
    is_current BOOLEAN DEFAULT false,
    created_by_id BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    published_at TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE(school_id, code)
);

CREATE INDEX idx_schedules_school ON schedules(school_id);
CREATE INDEX idx_schedules_status ON schedules(status);
CREATE INDEX idx_schedules_current ON schedules(is_current);

-- ============================================================================
-- SCHEDULE ENTRIES (Actual Scheduled Classes)
-- ============================================================================

CREATE TABLE schedule_entries (
    id BIGINT GENERATED ALWAYS AS IDENTITY (START WITH 1000 INCREMENT BY 1) PRIMARY KEY,
    schedule_id BIGINT NOT NULL REFERENCES schedules(id) ON DELETE CASCADE,
    course_id BIGINT NOT NULL REFERENCES courses(id) ON DELETE RESTRICT,
    teacher_id BIGINT NOT NULL REFERENCES teachers(id) ON DELETE RESTRICT,
    classroom_id BIGINT NOT NULL REFERENCES classrooms(id) ON DELETE RESTRICT,
    timeslot_id BIGINT NOT NULL REFERENCES timeslots(id) ON DELETE RESTRICT,
    
    -- Denormalized for easier access
    day_of_week VARCHAR(20) NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    course_code VARCHAR(50) NOT NULL,
    course_name VARCHAR(255) NOT NULL,
    teacher_name VARCHAR(200) NOT NULL,
    classroom_name VARCHAR(100) NOT NULL,
    
    -- Tracking
    enrollment_count INT DEFAULT 0,
    is_locked BOOLEAN DEFAULT false, -- Cannot be changed after lock
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_schedule_entries_schedule ON schedule_entries(schedule_id);
CREATE INDEX idx_schedule_entries_course ON schedule_entries(course_id);
CREATE INDEX idx_schedule_entries_teacher ON schedule_entries(teacher_id);
CREATE INDEX idx_schedule_entries_classroom ON schedule_entries(classroom_id);
CREATE INDEX idx_schedule_entries_timeslot ON schedule_entries(timeslot_id);
CREATE INDEX idx_schedule_entries_day_time ON schedule_entries(day_of_week, start_time);

-- ============================================================================
-- VIEWS FOR EASY QUERYING
-- ============================================================================

-- Schedule summary view
CREATE VIEW schedule_summary AS
SELECT 
    s.id,
    s.code as schedule_code,
    s.name as schedule_name,
    s.status,
    COUNT(DISTINCT se.id) as total_entries,
    COUNT(DISTINCT se.course_id) as unique_courses,
    COUNT(DISTINCT se.teacher_id) as unique_teachers,
    COUNT(DISTINCT se.classroom_id) as unique_classrooms,
    sch.name as school_name
FROM schedules s
JOIN schools sch ON s.school_id = sch.id
LEFT JOIN schedule_entries se ON s.id = se.schedule_id
GROUP BY s.id, s.code, s.name, s.status, sch.name;

-- Teacher schedule view
CREATE VIEW teacher_schedule AS
SELECT 
    t.id as teacher_id,
    t.first_name || ' ' || t.last_name as teacher_name,
    se.day_of_week,
    se.start_time,
    se.end_time,
    se.course_name,
    se.classroom_name,
    s.code as schedule_code
FROM schedule_entries se
JOIN teachers t ON se.teacher_id = t.id
JOIN schedules s ON se.schedule_id = s.id
ORDER BY se.day_of_week, se.start_time;

-- Classroom usage view
CREATE VIEW classroom_usage AS
SELECT 
    c.id as classroom_id,
    c.name as classroom_name,
    c.capacity,
    se.day_of_week,
    se.start_time,
    se.end_time,
    se.course_name,
    t.first_name || ' ' || t.last_name as teacher_name,
    s.code as schedule_code
FROM schedule_entries se
JOIN classrooms c ON se.classroom_id = c.id
JOIN teachers t ON se.teacher_id = t.id
JOIN schedules s ON se.schedule_id = s.id
ORDER BY c.name, se.day_of_week, se.start_time;

-- Course schedule view
CREATE VIEW course_schedule AS
SELECT 
    c.id as course_id,
    c.code as course_code,
    c.name as course_name,
    t.first_name || ' ' || t.last_name as teacher_name,
    se.day_of_week,
    se.start_time,
    se.end_time,
    cl.name as classroom_name,
    s.code as schedule_code
FROM schedule_entries se
JOIN courses c ON se.course_id = c.id
JOIN teachers t ON se.teacher_id = t.id
JOIN classrooms cl ON se.classroom_id = cl.id
JOIN schedules s ON se.schedule_id = s.id
ORDER BY c.code, se.day_of_week, se.start_time;

-- ============================================================================
-- ANALYTICAL VIEWS
-- ============================================================================

-- Teacher workload view
CREATE VIEW teacher_workload AS
SELECT 
    t.id as teacher_id,
    t.first_name || ' ' || t.last_name as teacher_name,
    s.code as schedule_code,
    COUNT(DISTINCT se.id) as total_classes,
    COUNT(DISTINCT se.day_of_week) as days_per_week,
    SUM((EXTRACT(HOUR FROM se.end_time) - EXTRACT(HOUR FROM se.start_time))
        + (EXTRACT(MINUTE FROM se.end_time) - EXTRACT(MINUTE FROM se.start_time))/60.0) as total_hours,
    MAX(EXTRACT(HOUR FROM se.end_time) - EXTRACT(HOUR FROM se.start_time)) as max_consecutive_hours
FROM schedule_entries se
JOIN teachers t ON se.teacher_id = t.id
JOIN schedules s ON se.schedule_id = s.id
GROUP BY t.id, t.first_name, t.last_name, s.code;

-- Classroom occupancy view
CREATE VIEW classroom_occupancy AS
SELECT 
    c.id as classroom_id,
    c.name as classroom_name,
    c.capacity,
    s.code as schedule_code,
    COUNT(DISTINCT se.id) as total_classes,
    COUNT(DISTINCT se.day_of_week) as days_used,
    COUNT(DISTINCT se.teacher_id) as unique_teachers,
    ROUND(100.0 * COUNT(DISTINCT se.id) / 
        (SELECT COUNT(DISTINCT timeslots.id) 
         FROM timeslots 
         WHERE timeslots.school_id = c.school_id), 2) as utilization_rate
FROM schedule_entries se
JOIN classrooms c ON se.classroom_id = c.id
JOIN schedules s ON se.schedule_id = s.id
GROUP BY c.id, c.name, c.capacity, s.code;

-- Scheduling conflicts view (to find issues)
CREATE VIEW potential_conflicts AS
SELECT 
    se1.id as entry1_id,
    se1.course_name as course1,
    se1.teacher_name as teacher1,
    se2.id as entry2_id,
    se2.course_name as course2,
    se2.teacher_name as teacher2,
    se1.day_of_week,
    se1.start_time,
    se1.end_time,
    CASE 
        WHEN se1.teacher_id = se2.teacher_id THEN 'Teacher Conflict'
        WHEN se1.classroom_id = se2.classroom_id THEN 'Room Conflict'
        ELSE 'Unknown'
    END as conflict_type
FROM schedule_entries se1
JOIN schedule_entries se2 ON 
    se1.schedule_id = se2.schedule_id AND
    se1.id < se2.id AND
    se1.day_of_week = se2.day_of_week AND
    se1.start_time = se2.start_time AND
    (se1.teacher_id = se2.teacher_id OR se1.classroom_id = se2.classroom_id)
ORDER BY se1.day_of_week, se1.start_time;

-- ============================================================================
-- CONSTRAINTS AND VALIDATIONS
-- ============================================================================

-- Ensure start_time is before end_time for timeslots
ALTER TABLE timeslots ADD CONSTRAINT check_timeslot_times
    CHECK (start_time < end_time);

-- Ensure schedule dates are valid
ALTER TABLE schedules ADD CONSTRAINT check_schedule_dates
    CHECK (start_date < end_date);

-- Ensure course dates are within schedule dates
ALTER TABLE schedule_entries ADD CONSTRAINT check_entry_dates
    CHECK (start_time < end_time);

-- ============================================================================
-- SAMPLE DATA (Optional - for testing)
-- ============================================================================

-- Insert sample school
INSERT INTO schools (name, code, city, state, principal_name, is_active)
VALUES ('Central University', 'CU-001', 'New York', 'NY', 'Dr. James Smith', true);

-- Insert sample teachers
INSERT INTO teachers (school_id, employee_id, first_name, last_name, email, specialization, hire_date, is_active)
VALUES 
((SELECT id FROM schools WHERE code = 'CU-001'), '1001', 'John', 'Smith', 'john@university.edu', 'Physics', '2020-01-15', true),
((SELECT id FROM schools WHERE code = 'CU-001'), '1002', 'Emily', 'Davis', 'emily@university.edu', 'Mathematics', '2019-06-20', true),
((SELECT id FROM schools WHERE code = 'CU-001'), '1003', 'Michael', 'Johnson', 'michael@university.edu', 'Chemistry', '2021-08-10', true),
((SELECT id FROM schools WHERE code = 'CU-001'), '1004', 'Sarah', 'Williams', 'sarah@university.edu', 'Biology', '2018-03-05', true);

-- Insert sample classrooms
INSERT INTO classrooms (school_id, code, name, room_type, capacity, floor_number, building, is_active)
VALUES 
((SELECT id FROM schools WHERE code = 'CU-001'), 'LAB-A', 'Science Lab A', 'Lab', 25, 2, 'Building A', true),
((SELECT id FROM schools WHERE code = 'CU-001'), 'LAB-B', 'Science Lab B', 'Lab', 25, 2, 'Building B', true),
((SELECT id FROM schools WHERE code = 'CU-001'), 'CR-101', 'Classroom 101', 'Classroom', 40, 1, 'Building A', true),
((SELECT id FROM schools WHERE code = 'CU-001'), 'CR-102', 'Classroom 102', 'Classroom', 40, 1, 'Building A', true),
((SELECT id FROM schools WHERE code = 'CU-001'), 'LH-1', 'Lecture Hall 1', 'Auditorium', 100, 1, 'Building B', true);

-- Insert sample timeslots
INSERT INTO timeslots (school_id, day_of_week, start_time, end_time, slot_order, is_active)
VALUES 
((SELECT id FROM schools WHERE code = 'CU-001'), 'Monday', '09:00', '10:00', 1, true),
((SELECT id FROM schools WHERE code = 'CU-001'), 'Monday', '10:00', '11:00', 2, true),
((SELECT id FROM schools WHERE code = 'CU-001'), 'Monday', '11:00', '12:00', 3, true),
((SELECT id FROM schools WHERE code = 'CU-001'), 'Monday', '13:00', '14:00', 4, true),
((SELECT id FROM schools WHERE code = 'CU-001'), 'Tuesday', '09:00', '10:00', 1, true),
((SELECT id FROM schools WHERE code = 'CU-001'), 'Tuesday', '10:00', '11:00', 2, true),
((SELECT id FROM schools WHERE code = 'CU-001'), 'Wednesday', '09:00', '10:00', 1, true),
((SELECT id FROM schools WHERE code = 'CU-001'), 'Thursday', '09:00', '10:00', 1, true),
((SELECT id FROM schools WHERE code = 'CU-001'), 'Friday', '09:00', '10:00', 1, true);

-- Insert sample courses
INSERT INTO courses (school_id, code, name, teacher_id, semester, duration_minutes, required_days_per_week, min_day_gap, preferred_days, is_active)
VALUES 
((SELECT id FROM schools WHERE code = 'CU-001'), 'PHY-101', 'Physics I', 
 (SELECT id FROM teachers WHERE employee_id = '1001'), 1, 60, 3, 1, 'Monday,Wednesday,Friday', true),
((SELECT id FROM schools WHERE code = 'CU-001'), 'MATH-101', 'Calculus I', 
 (SELECT id FROM teachers WHERE employee_id = '1002'), 1, 60, 3, 1, 'Monday,Wednesday,Friday', true),
((SELECT id FROM schools WHERE code = 'CU-001'), 'CHM-101', 'Chemistry I', 
 (SELECT id FROM teachers WHERE employee_id = '1003'), 1, 60, 2, 2, 'Tuesday,Thursday', true),
((SELECT id FROM schools WHERE code = 'CU-001'), 'BIO-101', 'Biology I', 
 (SELECT id FROM teachers WHERE employee_id = '1004'), 1, 60, 2, 2, 'Tuesday,Thursday', true);

-- Insert sample schedule
INSERT INTO schedules (school_id, name, code, start_date, end_date, status, is_current, created_at)
VALUES 
((SELECT id FROM schools WHERE code = 'CU-001'), 'Fall 2024', 'FALL-2024', '2024-09-01', '2024-12-15', 'Published', true, CURRENT_TIMESTAMP);

-- Insert sample schedule entries
INSERT INTO schedule_entries (
    schedule_id, course_id, teacher_id, classroom_id, timeslot_id,
    day_of_week, start_time, end_time, course_code, course_name, teacher_name, classroom_name
)
VALUES 
(
    (SELECT id FROM schedules WHERE code = 'FALL-2024'),
    (SELECT id FROM courses WHERE code = 'PHY-101'),
    (SELECT id FROM teachers WHERE employee_id = '1001'),
    (SELECT id FROM classrooms WHERE code = 'LAB-A'),
    (SELECT id FROM timeslots WHERE day_of_week = 'Monday' AND start_time = '09:00'),
    'Monday', '09:00', '10:00', 'PHY-101', 'Physics I', 'John Smith', 'Science Lab A'
),
(
    (SELECT id FROM schedules WHERE code = 'FALL-2024'),
    (SELECT id FROM courses WHERE code = 'MATH-101'),
    (SELECT id FROM teachers WHERE employee_id = '1002'),
    (SELECT id FROM classrooms WHERE code = 'CR-101'),
    (SELECT id FROM timeslots WHERE day_of_week = 'Monday' AND start_time = '10:00'),
    'Monday', '10:00', '11:00', 'MATH-101', 'Calculus I', 'Emily Davis', 'Classroom 101'
),
(
    (SELECT id FROM schedules WHERE code = 'FALL-2024'),
    (SELECT id FROM courses WHERE code = 'CHM-101'),
    (SELECT id FROM teachers WHERE employee_id = '1003'),
    (SELECT id FROM classrooms WHERE code = 'LAB-B'),
    (SELECT id FROM timeslots WHERE day_of_week = 'Tuesday' AND start_time = '09:00'),
    'Tuesday', '09:00', '10:00', 'CHM-101', 'Chemistry I', 'Michael Johnson', 'Science Lab B'
);

-- ============================================================================
-- END OF SCHEMA
-- ============================================================================

COMMIT;

-- Verification
-- SELECT 'Schools:' as entity, COUNT(*) as count FROM schools
-- UNION ALL
-- SELECT 'Teachers:', COUNT(*) FROM teachers
-- UNION ALL
-- SELECT 'Students:', COUNT(*) FROM students
-- UNION ALL
-- SELECT 'Courses:', COUNT(*) FROM courses
-- UNION ALL
-- SELECT 'Classrooms:', COUNT(*) FROM classrooms
-- UNION ALL
-- SELECT 'Timeslots:', COUNT(*) FROM timeslots
-- UNION ALL
-- SELECT 'Constraints:', COUNT(*) FROM constraints
-- UNION ALL
-- SELECT 'Schedules:', COUNT(*) FROM schedules
-- UNION ALL
-- SELECT 'Schedule Entries:', COUNT(*) FROM schedule_entries;
