### Location
___
- Sits in between and communicates to Lock Manager and Notification Micro Service. 
#### Purpose
____
- The **Session Manager** handles the **lifecycle and coordination** of a Guardian–User session.
#### Responsibilities 
____
- Create, update, and end **sessions** (e.g., `startSession(userId, guardianId)`).
    
- Associate `tenant_id`, `session_id`, `user_id`, and `guardian_id`.
    
- Store **session metadata**:
    
    - Start time, end time
        
    - Whether a session is active
        
    - Optional session-specific rules (e.g., "study mode")
        
- Validate if a session is still active for other services (e.g., Lock Manager, Notification)

#### Outputs
___
- Session state (`active`, `expired`, `revoked`)
    
- Session config (references to restrictions, lock presets, etc.)


#### Tables 

```postgres
-- Tenants represent organizations or groups using the system
CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);

-- Guardians are the controllers/parents/accountability partners
CREATE TABLE guardians (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    email TEXT NOT NULL UNIQUE,
    name TEXT,
    created_at TIMESTAMPTZ DEFAULT now()
);

-- Users are the people being monitored
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    guardian_id UUID NOT NULL REFERENCES guardians(id) ON DELETE CASCADE,
    email TEXT NOT NULL UNIQUE,
    name TEXT,
    created_at TIMESTAMPTZ DEFAULT now()
);

-- A session links one user and one guardian and governs an enforcement period
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    guardian_id UUID NOT NULL REFERENCES guardians(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status TEXT NOT NULL CHECK (status IN ('active', 'ended', 'expired')),
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Optional table to hold custom settings for a session (lock policies, labels, etc.)
CREATE TABLE session_configs (
    session_id UUID PRIMARY KEY REFERENCES sessions(id) ON DELETE CASCADE,
    label TEXT,
    lock_profile JSONB, -- e.g. {"lockedApps": ["tiktok", "youtube"], "whitelist": ["duolingo"]}
    notify_on_unlock BOOLEAN DEFAULT true,
    custom_message TEXT,
    created_at TIMESTAMPTZ DEFAULT now()
);
```
