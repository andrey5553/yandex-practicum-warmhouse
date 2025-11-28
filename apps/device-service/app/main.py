from fastapi import FastAPI
from .routes import router
from .database import create_tables
from .config import settings

app = FastAPI(
    title=settings.app_name,
    version=settings.version,
    description="Микросервис для управления IoT устройствами умного дома"
)

# Include routes
app.include_router(router, prefix="/api/v1")

# Create database tables on startup
@app.on_event("startup")
def on_startup():
    create_tables()

@app.get("/")
def read_root():
    return {"message": "Device Service API", "version": settings.version}

@app.get("/health")
def health_check():
    return {"status": "healthy"}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host=settings.host, port=settings.port)
