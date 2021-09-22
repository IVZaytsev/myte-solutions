from datetime import date
from django.db import models
from django.core.validators import MaxValueValidator, MinValueValidator

# Create your models here.
class Position(models.Model):
    desc = models.CharField(max_length=30, verbose_name='Название')
    salary = models.IntegerField(default=15000, validators=[MaxValueValidator(1000000), MinValueValidator(12792)], verbose_name='Зар.плата')    # МРОТ

    def __str__(self):
        return self.desc
        
    class Meta:
        verbose_name = 'Должность'
        verbose_name_plural = 'Должности'

class Worker(models.Model):
    fullname = models.CharField(max_length=60, verbose_name='ФИО')
    dateEmploy = models.DateField(default=date.today, validators=[MaxValueValidator(limit_value=date.today)], verbose_name='Принят')            # auto_now_add=True
    position = models.ForeignKey(Position, null=True, blank=True, on_delete=models.SET_NULL, verbose_name='Должность')

    def __str__(self):
        return self.fullname
        
    class Meta:
        verbose_name = 'Работник'
        verbose_name_plural = 'Работники'

class Department(models.Model):
    desc = models.CharField(max_length=30, verbose_name='Название')
    chief = models.ForeignKey(Worker, null=True, blank=True, on_delete=models.SET_NULL, related_name='chief', verbose_name='Руководитель')
    workers = models.ManyToManyField(Worker, blank=True, related_name='workers', verbose_name='Работники')

    def __str__(self):
        return self.desc

    def get_worker_count(self):
        return len(self.workers.all())

    class Meta:
        verbose_name = 'Отдел'
        verbose_name_plural = 'Отделы'
