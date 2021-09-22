# Generated by Django 3.2.7 on 2021-09-22 13:32

import datetime
import django.core.validators
from django.db import migrations, models
import django.db.models.deletion


class Migration(migrations.Migration):

    initial = True

    dependencies = [
    ]

    operations = [
        migrations.CreateModel(
            name='Position',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('desc', models.CharField(max_length=30, verbose_name='Название')),
                ('salary', models.IntegerField(default=15000, validators=[django.core.validators.MaxValueValidator(1000000), django.core.validators.MinValueValidator(12792)], verbose_name='Зар.плата')),
            ],
            options={
                'verbose_name': 'Должность',
                'verbose_name_plural': 'Должности',
            },
        ),
        migrations.CreateModel(
            name='Worker',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('fullname', models.CharField(max_length=60, verbose_name='ФИО')),
                ('dateEmploy', models.DateField(default=datetime.date.today, validators=[django.core.validators.MaxValueValidator(limit_value=datetime.date.today)], verbose_name='Принят')),
                ('position', models.ForeignKey(blank=True, null=True, on_delete=django.db.models.deletion.SET_NULL, to='workers.position', verbose_name='Должность')),
            ],
            options={
                'verbose_name': 'Работник',
                'verbose_name_plural': 'Работники',
            },
        ),
        migrations.CreateModel(
            name='Department',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('desc', models.CharField(max_length=30, verbose_name='Название')),
                ('chief', models.ForeignKey(blank=True, null=True, on_delete=django.db.models.deletion.SET_NULL, related_name='chief', to='workers.worker', verbose_name='Руководитель')),
                ('workers', models.ManyToManyField(blank=True, related_name='workers', to='workers.Worker', verbose_name='Работники')),
            ],
            options={
                'verbose_name': 'Отдел',
                'verbose_name_plural': 'Отделы',
            },
        ),
    ]
