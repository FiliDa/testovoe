-- Add periodicity settings to tasks table
ALTER TABLE tasks 
ADD COLUMN periodicity_type TEXT,
ADD COLUMN periodicity_interval INTEGER,
ADD COLUMN periodicity_days_of_month INTEGER[],
ADD COLUMN periodicity_specific_dates DATE[],
ADD COLUMN periodicity_even_odd TEXT,
ADD COLUMN next_occurrence DATE;

-- Create index for efficient querying of recurring tasks
CREATE INDEX IF NOT EXISTS idx_tasks_next_occurrence ON tasks (next_occurrence) 
WHERE periodicity_type IS NOT NULL;