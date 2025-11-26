import os
from pydantic_settings import BaseSettings

class Settings(BaseSettings):
    app_name: str = "Device Service"
    version: str = "1.0.0"
    
    # Database
    database_url: str = os.getenv(
        "DATABASE_URL", 
        "postgresql://user:password@localhost:5433/device_service"
    )
    
    # Service
    host: str = os.getenv("HOST", "0.0.0.0")
    port: int = int(os.getenv("PORT", 8083))
    
    class Config:
        env_file = ".env"

settings = Settings()
